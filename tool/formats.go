package tool

//func (f *FileFormat) Mimetype() string {
//  return mime.TypeByExtension(f.Ext())
//}

//func (f FileFormat) IsImage() bool {
//  if strings.Contains(f.Mimetype(), "image") {
//    return true
//  }
//  return false
//}

//func (f FileFormat) IsAudio() bool {
//  if strings.Contains(f.Mimetype(), "audio") {
//    return true
//  } else {
//    fmt.Println("not an audio file")
//  }
//  return false
//}

//func (f FileFormat) IsPlainText() bool {
//  if strings.Contains(f.Mimetype(), "text/plain") {
//    return true
//  } else {
//    log.Fatalln("needs to be plain text file")
//  }
//  return false
//}

//func (f FileFormat) IsFFmeta() bool {
//  if f.IsPlainText() {
//    contents, err := os.Open(f.Path())
//    if err != nil {
//      log.Fatal(err)
//    }
//    defer contents.Close()

//    scanner := bufio.NewScanner(contents)
//    line := 0
//    for scanner.Scan() {
//      if line == 0 && scanner.Text() == ";FFMETADATA1" {
//        return true
//      } else {
//        log.Fatalln("ffmpeg metadata files need to have ';FFMETADATA1' as the first line")
//      }
//    }
//  }
//  return false
//}
