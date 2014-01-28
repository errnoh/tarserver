// Copyright 2014 errnoh. All rights reserved.
// Use of this source code is governed by
// MIT License that can be found in the LICENSE file.
package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io"
	"log"
	"net/http"
	"time"
)

func init() {
	supported["tar"] = tarrer
	supported["gz"] = gzipper
	supported["gzip"] = gzipper
	supported["zlib"] = zlibber
	supported["echo"] = echo
}

func tarrer(req *http.Request, r io.Reader, w *bytes.Buffer) {
	w.Reset()

	tw := tar.NewWriter(w)
	mpr, err := req.MultipartReader()
	if err != nil {
		log.Println(err)
		return // XXX: return empty if not a multipart?
	}

	var b bytes.Buffer
	for f, err := mpr.NextPart(); ; f, err = mpr.NextPart() {
		if err != nil {
			log.Println(err)
			break
		}
		defer f.Close()
		b.Reset()
		n, err := b.ReadFrom(f)
		if err != nil {
			break
		}

		hdr := &tar.Header{
			Name:    f.FileName(),
			ModTime: time.Now(),
			Size:    n,
		}
		if err = tw.WriteHeader(hdr); err != nil {
			log.Println(err)
			w.Reset()
			break
		}
		if _, err := io.Copy(tw, &b); err != nil {
			log.Println(err)
			w.Reset()
			break
		}
	}
	if err = tw.Close(); err != nil {
		log.Println(err)
		w.Reset()
	}
}

func gzipper(req *http.Request, r io.Reader, w *bytes.Buffer) {
	w.Reset()
	gzw := gzip.NewWriter(w)
	io.Copy(gzw, r)
	gzw.Flush()
	gzw.Close()
}

func zlibber(req *http.Request, r io.Reader, w *bytes.Buffer) {
	w.Reset()
	zw := zlib.NewWriter(w)
	io.Copy(zw, r)
	zw.Flush()
	zw.Close()
}

func echo(req *http.Request, r io.Reader, w *bytes.Buffer) {
	w.Reset()
	io.Copy(w, r)
}
