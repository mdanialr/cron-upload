package s3

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/mdanialr/cron-upload/internal/provider"
	h "github.com/mdanialr/cron-upload/pkg/helper"
)

// NewS3BucketProvider return provider that use AWS S3 Bucket as the cloud
// provider.
func NewS3BucketProvider(ctx context.Context, bucket string, svc *s3.Client) provider.Cloud {
	return &s3Bucket{
		ctx:    ctx,
		bucket: bucket,
		svc:    svc,
	}
}

type s3Bucket struct {
	ctx    context.Context
	bucket string
	svc    *s3.Client
}

func (s *s3Bucket) GetFolders(parent ...string) ([]*provider.Payload, error) {
	objsInput := &s3.ListObjectsV2Input{
		Bucket: h.Ptr(s.bucket),
	}
	// additionally use the parent as the prefix if provided
	if len(parent) > 0 {
		parentKey := strings.TrimSuffix(parent[0], "/") + "/"
		parentKey = strings.TrimPrefix(parentKey, "/") // make sure has no leading slash
		objsInput.StartAfter = h.Ptr(parentKey)
	}
	objs, err := s.svc.ListObjectsV2(s.ctx, objsInput)
	if err != nil {
		return nil, fmt.Errorf("failed to query for folders: %s", err)
	}
	if objs == nil {
		return nil, fmt.Errorf("no folder was found")
	}
	// transform to local payload type
	var res []*provider.Payload
	for _, obj := range objs.Contents {
		objKey := h.Def(obj.Key)
		if strings.HasSuffix(objKey, "/") {
			newPayload := &provider.Payload{
				Name:   strings.TrimSuffix(h.Def(obj.Key), "/"),
				Id:     strings.TrimSuffix(h.Def(obj.Key), "/"),
				Parent: parent,
			}
			if obj.LastModified != nil {
				newPayload.CreatedAt = obj.LastModified.Format(time.RFC3339)
			}
			res = append(res, newPayload)
		}
	}
	return res, nil
}

func (s *s3Bucket) CreateFolder(name string, parent ...string) (string, error) {
	name = strings.TrimSuffix(name, "/") + "/" // make sure to manually add trailing slash
	obj := &s3.PutObjectInput{
		Bucket: h.Ptr(s.bucket),
		Key:    h.Ptr(name),
		ACL:    types.ObjectCannedACLPrivate,
	}
	// additionally add parent if provided
	var parentKey string
	if len(parent) > 0 {
		parentKey = strings.TrimSuffix(parent[0], "/") + "/"
		parentKey = strings.TrimPrefix(parentKey, "/") // make sure has no leading slash
		obj.Key = h.Ptr(parentKey + name)
	}
	if _, err := s.svc.PutObject(s.ctx, obj); err != nil {
		return "", fmt.Errorf("failed to create a folder with name '%s': %s", name, err)
	}
	// return along with the parent if any
	if parentKey != "" {
		return parentKey + name, nil
	}
	return name, nil
}

func (s *s3Bucket) GetFiles(folderId string) ([]*provider.Payload, error) {
	parentKey := strings.TrimSuffix(folderId, "/") + "/"
	parentKey = strings.TrimPrefix(parentKey, "/") // make sure has no leading slash
	objsInput := &s3.ListObjectsV2Input{
		Bucket:     h.Ptr(s.bucket),
		StartAfter: h.Ptr(parentKey),
	}
	objs, err := s.svc.ListObjectsV2(s.ctx, objsInput)
	if err != nil {
		return nil, fmt.Errorf("failed to query for folders: %s", err)
	}
	if objs != nil {
		if len(objs.Contents) > 0 {
			var res []*provider.Payload
			for _, obj := range objs.Contents {
				newPayload := &provider.Payload{
					Name:   strings.TrimSuffix(h.Def(obj.Key), "/"),
					Id:     strings.TrimSuffix(h.Def(obj.Key), "/"),
					Parent: []string{parentKey},
				}
				if obj.LastModified != nil {
					newPayload.CreatedAt = obj.LastModified.Format(time.RFC3339)
				}
				res = append(res, newPayload)
			}
			return res, nil
		}
	}
	return nil, nil
}

func (s *s3Bucket) UploadFile(payload *provider.Payload, chunkSize ...int) (*provider.Payload, error) {
	defer payload.File.Close()

	// use manager to upload with fixed file size
	uploader := manager.NewUploader(s.svc, func(u *manager.Uploader) {
		if len(chunkSize) > 0 {
			u.PartSize = int64(chunkSize[0] * 1024)
		}
	})
	// additionally append the parent if provided
	if len(payload.Parent) > 0 {
		parentKey := strings.TrimSuffix(payload.Parent[0], "/") + "/"
		parentKey = strings.TrimPrefix(parentKey, "/") // make sure has no leading slash
		payload.Name = parentKey + payload.Name
	}
	obj := &s3.PutObjectInput{
		Bucket: h.Ptr(s.bucket),
		Key:    h.Ptr(payload.Name),
		Body:   payload.File,
		ACL:    types.ObjectCannedACLPrivate,
	}
	if _, err := uploader.Upload(s.ctx, obj); err != nil {
		return nil, fmt.Errorf("failed to upload object with key '%s': %s", payload.Name, err)
	}

	resPayload := &provider.Payload{
		Name:   strings.TrimSuffix(payload.Name, "/"),
		Id:     strings.TrimSuffix(payload.Name, "/"),
		Parent: payload.Parent,
	}
	return resPayload, nil
}

func (s *s3Bucket) Delete(id string) error {
	obj := &s3.DeleteObjectInput{
		Bucket: h.Ptr(s.bucket),
		Key:    h.Ptr(id),
	}
	if _, err := s.svc.DeleteObject(s.ctx, obj); err != nil {
		return fmt.Errorf("failed to delete an object with key '%s': %s", id, err)
	}
	return nil
}
