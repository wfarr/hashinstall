package hashinstall

import (
	"archive/zip"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
)

// Install will fetch the Hashicorp package of name and version and install it
func Install(name, version, destdir string, info chan string, debug chan string) error {

	// fetch the zipball from hashicorp releases
	url := fmt.Sprintf("https://releases.hashicorp.com/%v/%v/%v_%v_%v_%v.zip", name, version, name, version, runtime.GOOS, runtime.GOARCH)

	info <- fmt.Sprintln("downloading archive from", url)

	resp, err := http.Get(url)
	if err != nil {
		if resp.StatusCode == http.StatusForbidden {
			return errors.New("could not download release: you may have the wrong name or version")
		}

		return errors.New("could not download release")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// write the zipball to a temp file
	tmpfile, err := ioutil.TempFile("", fmt.Sprintf("%v_%v_%v_%v.zip", name, version, runtime.GOOS, runtime.GOARCH))
	if err != nil {
		return err
	}
	defer tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(body); err != nil {
		return err
	}

	// read the zipball
	r, err := zip.OpenReader(tmpfile.Name())
	if err != nil {
		return err
	}
	defer r.Close()

	// ensure we can write to destdir
	err = os.MkdirAll(destdir, 0755)
	if err != nil {
		return err
	}

	// for each file in the zipball, write it out to a file of the same name inside the destdir
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}

		content, err := ioutil.ReadAll(rc)
		if err != nil {
			return err
		}

		// if f.Name is inside a subdirectory in the zipball, we need to make a directory to put it in
		i := strings.LastIndexByte(f.Name, '/')
		if i > -1 {
			err = os.MkdirAll(fmt.Sprintf("%v/%v", destdir, f.Name[:i]), 0755)
			if err != nil {
				return err
			}
		}

		// write archived file out to destdir
		err = ioutil.WriteFile(fmt.Sprintf("%v/%v", destdir, f.Name), content, 0755)
		if err != nil {
			return err
		}
		rc.Close()
	}

	if err := tmpfile.Close(); err != nil {
		return err
	}

	return nil
}
