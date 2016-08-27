package   cmd

import (
  "path/filepath"
  "github.com/spf13/viper"
  "io/ioutil"
)

func config(key string) string {
  return viper.GetString(key)
}


type TerraformLayer struct {
  RequiredVars []string
  SourceUri string
  Environment string
  Name string
}

func (t *TerraformLayer) dir() string {
  dirPath := filepath.Join(config("terraform-dir"),config("env"),t.Name)
  err := os.MkdirAll(dirPath, dirMode)
  if err != nil {
    log.Fatalf("Unable to create directory at %s:%s", dirPath, err)
  }
  return dirPath
}


func (t *TerraformLayer) makeMake() bool {

  dir := t.dir()
  finfo, err := os.Stat(filepath.Join(dir,"Makefile"))
  if err != nil  || OverWriteFiles {
    makeContents,err := Asset(filepath.Join("assets", "Makefile"))
		if err != nil {
			log.Fatalln("Unable to retrieve base Makefile file", err)
		}
    err = ioutil.WriteFile("/tmp/dat1", d1, 0644)
    defer finfo.Close()

    return true
  }
  return false
}



func makeMake(path String) {
  finfo, err := os.Stat(filepath.Join(path,"Makefile"))
  if err != nil {

    // no such file or dir
    return
  }
  if finfo.IsDir() {
    // it's a file
  } else {
    // it's a directory
  }



}

func terraformSkeleton(path string) {
  templ := ParseTemplate("s3terraform", "s3.tf")
  f := OpenForWrite("bucket","main.tf")
  if err := templ.Execute(f, bucketRequest); err != nil {
    log.Fatalln("Unable to generate template", err)
  }
}




func OpenForWrite(modName,fileName string) *os.File {
  dirPath := TerraformPath(modName)
  err := os.MkdirAll(dirPath, dirMode)
  if err != nil {
		log.Fatalf("Unable to create directory at %s:%s", dirPath, err)
	}

	file := filepath.Join(dirPath, fileName)
	f, err := os.Create(file)
  if err != nil {
		log.Fatalf("Unable to create file at %s:%v", file, err)
	}
  return f
}
