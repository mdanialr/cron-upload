package gdrive

// LISTING ALL FILES and FOLDERS within given fileId (which is the id of a folder).
//list, _ := dr.Files.List().Q("'1Zl6s95uByQGAplNbtjVgbeXeE6x48vqL' in parents").Do()
//for _, fl := range list.Files {
//	fmt.Println("Id:", fl.Id)
//
//	b, _ := fl.MarshalJSON()
//	var out bytes.Buffer
//	json.Indent(&out, b, "", "\t")
//	fmt.Println(out.String() + "\n")
//}

// LISTING ALL FILES (only) within a folder from the given id. retrieve only non-folder data.
//const id = "1I0Ou4Uu0S_RhIdR2RYjPXiBa2TyppLzn"
//query := fmt.Sprintf("'%s' in parents and %s != '%s'", id, gdrive.FieldMIME, gdrive.MIMEFolder)
//list, _ := dr.Files.List().Q(query).Do()
//for _, fl := range list.Files {
//	b, _ := fl.MarshalJSON()
//	var out bytes.Buffer
//	json.Indent(&out, b, "", "\t")
//	fmt.Println(out.String() + "\n")
//}

// DELETE a FOLDER
//const id = "1mdGZuCobVn6grxZJmNzbPeXIS9Jk0ztH"
//if err := dr.Files.Delete(id).Do(); err != nil {
//	log.Fatalf("failed to delete a file or a folder with id: %s and error: %s \n", id, err)
//}

// CREATE a FOLDER
//fl := &drive.File{Name: "VPS-Backup", MimeType: gdrive.MIMEFolder}
//newFl, err := dr.Files.Create(fl).Do()
//if err != nil {
//	log.Fatalln("failed to create a folder:", err)
//}
//b, _ := newFl.MarshalJSON()
//var out bytes.Buffer
//json.Indent(&out, b, "", "\t")
//fmt.Println(out.String() + "\n")

// GET a FOLDER with a NAME of 'VPS-Backup'
//q := fmt.Sprintf("mimeType = '%s' and name = '%s'", gdrive.MIMEFolder, "VPS-Backup")
//folder, err := dr.Files.List().Q(q).Do()
//if err != nil {
//	log.Fatalln("failed to query for a folder with a name VPS-Backup:", err)
//}
//// List the result only if the length is not zero
//if len(folder.Files) > 0 {
//	b, _ := folder.Files[0].MarshalJSON()
//	var out bytes.Buffer
//	json.Indent(&out, b, "", "\t")
//	fmt.Println(out.String() + "\n")
//	fmt.Println("ID:", folder.Files[0].Id)
//}

// QUERY using GIVEN MIME TYPE & has GIVEN String name
//q := fmt.Sprintf("%s = '%s' and name contains '%s'", gdrive.FieldMIME, gdrive.MIMEPlainText, "newInitServ")
//folder, err := dr.Files.List().Q(q).Do()
//if err != nil {
//	log.Fatalln("failed to query for a folder with a name VPS-Backup:", err)
//}
//// List the result only if the length is not zero
//if len(folder.Files) > 0 {
//	b, _ := folder.Files[0].MarshalJSON()
//	var out bytes.Buffer
//	json.Indent(&out, b, "", "\t")
//	fmt.Println(out.String() + "\n")
//	fmt.Println("ID:", folder.Files[0].Id)
//}

// RETRIEVE DETAIL OF THE FILE or FOLDER for GIVEN id
//const id = "1chTTpkL9n0Vsru_NdNl7fBxOuX1gwY9F"
//qFile, err := dr.Files.Get(id).Fields(
//	gdrive.FieldId,
//	gdrive.FieldMIME,
//	gdrive.FieldName,
//	gdrive.FieldParents,
//).Do()
//if err != nil {
//	log.Fatalf("failed to retrieve file with id: %s and error: %s\n", id, err)
//}
//b, _ := qFile.MarshalJSON()
//var out bytes.Buffer
//json.Indent(&out, b, "", "\t")
//fmt.Println(out.String() + "\n")
//fmt.Println("ID:", qFile.Id)

// UPLOAD a FILE IN THE DESIGNATED FOLDER which is VPS-Backup (1Cyo9qlbW7zzVTMJ7qk9E5n9ojwK--WT1)
//const filePath = "/home/user/ssh-port-forwarding.zip"
//flInstance, err := os.Open(filePath)
//if err != nil {
//	log.Fatalf("failed to open file: %s with error: %s\n", filePath, err)
//}
//defer flInstance.Close()
//fl := &drive.File{
//	Parents:  []string{"1Cyo9qlbW7zzVTMJ7qk9E5n9ojwK--WT1"},
//	MimeType: gdrive.MIMEZip, // this field is optional if the filename already has extension
//	Name:     filepath.Base(flInstance.Name()),
//}
//uploadFl, err := dr.Files.Create(fl).Media(flInstance).Fields(
//	gdrive.FieldId,
//	gdrive.FieldMIME,
//	gdrive.FieldName,
//	gdrive.FieldParents,
//).Do()
//if err != nil {
//	log.Fatalf("failed to upload file: %s with error: %s\n", uploadFl.Name, err)
//}
//b, _ := uploadFl.MarshalJSON()
//var out bytes.Buffer
//json.Indent(&out, b, "", "\t")
//fmt.Println(out.String() + "\n")
//fmt.Println("ID:", uploadFl.Id)

// EMPTY TRASH
//if err := dr.Files.EmptyTrash().Do(); err != nil {
//	log.Fatalln("failed tp emptying trash:", err)
//}
