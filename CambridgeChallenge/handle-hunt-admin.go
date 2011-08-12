/*
 respond-huntadmin.go - respond to requests for /huntadmin, an administrative 
 	      	      interface for adding, deleting, and modifying hunt data.

    Copyright 2011 Joshua C. Randall <jcrandall@alum.mit.edu>

    This file is part of CambridgeChallenge.

    CambridgeChallenge is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    CambridgeChallenge is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with CambridgeChallenge.  If not, see <http://www.gnu.org/licenses/>.
*/
package CambridgeChallenge

import (
    "appengine"
    "appengine/blobstore"
    "appengine/datastore"
    "http"
    "template"
    "os"
    "time"
    "regexp"
)

const (
      huntAdminPath = "/HuntAdmin/"
      huntAdminUploadPath = "/HuntAdmin/Upload/"
      huntAdminDownloadPath = "/HuntAdmin/Download/"
      huntAdminTemplateFileName = "template/huntadmin.html.gotmpl"
      huntDirectoryDatastore = "HuntDirectory"
      huntLimit = 100
)

type HuntDirectoryEntry struct {
    HuntName	   string
    BlobKey	   appengine.BlobKey
    CreatedDate	   datastore.Time
    Creator	   string
}

type HuntAdminTemplateData struct {
     User string
     UploadURL *http.URL
     UserError string
     StatusMessage string
     Hunts []HuntDirectoryEntry
}

func init() {
    http.HandleFunc(huntAdminPath, handleHuntAdmin)
    http.HandleFunc(huntAdminUploadPath, handleHuntAdminUpload)
    http.HandleFunc(huntAdminDownloadPath, handleHuntAdminDownload)
}

var huntadminTemplate = template.MustParseFile(huntAdminTemplateFileName, template.FormatterMap{"dstime" : dstimeFormatter})

var hdQuery = datastore.NewQuery(huntDirectoryDatastore).Order("-CreatedDate").Limit(huntLimit) //TODO this should ideally not be a hard limit


func handleHuntAdmin(w http.ResponseWriter, r *http.Request) {
    var td HuntAdminTemplateData
    var err os.Error
    td.User = requireAnyUser(w, r)
    LogAccess(r, td.User)

    c := appengine.NewContext(r)

    td.Hunts = make([]HuntDirectoryEntry, 0, huntLimit)
    if _, err := hdQuery.GetAll(c, &td.Hunts); err != nil {
       serveError(c, w, err)
       return
    }

    err = r.ParseForm()
    if err != nil {
      serveError(c, w, err)
      return
    }
    td.UserError = r.FormValue("user_error")
    td.StatusMessage = r.FormValue("status_message")

    td.UploadURL, err = blobstore.UploadURL(c, huntAdminUploadPath)
    if err != nil {
       serveError(c, w, err)
       return
    }

    w.Header().Set("Content-Type", "text/html")
    if err = huntadminTemplate.Execute(w, td); err != nil {
        serveError(c, w, err)
    }
}

    	      
func deleteBlobsOnRecover(w http.ResponseWriter, r *http.Request, blobs map[string][]*blobstore.BlobInfo) {
  if rec := recover(); rec != nil {
     c := appengine.NewContext(r)
     c.Infof("Recovering in deleteBlobsOnRecover from: %v", rec)
     // panic has occured, recover from it by deleting all blobs
     for _, blobInfos := range blobs {
     	 for _, blobInfo := range blobInfos {
	     err := blobstore.Delete(c, blobInfo.BlobKey)
	     if err != nil {
	     	c.Errorf("error deleting blobkey %v", blobInfo.BlobKey)
	     } else {
	        c.Debugf("blobkey %v deleted", blobInfo.BlobKey)
	     }
	 }
     }     
  }
  return
}



func handleHuntAdminUpload(w http.ResponseWriter, r *http.Request) {
     c := appengine.NewContext(r)
     user := requireAnyUser(w, r)
     blobs, other, err := blobstore.ParseUpload(r)
     if err != nil {
        serveError(c, w, err)
        return
     }
     // at this point, blob exists, so if we have an error we must delete it
     defer deleteBlobsOnRecover(w, r, blobs)

     if len(blobs) > 1 {
     	// more than one blob name uploaded, bail out
        http.Redirect(w, r, huntAdminPath+"/?user_error=multiple files uploaded, please use the form provided to upload a single JSON.", http.StatusFound)
	panic("multiple blobs uploaded")
     }

     huntJSONBlobInfos := blobs["hunt_json"] // must match form input element name
     if len(huntJSONBlobInfos) == 0 {
        c.Errorf("no hunt_json file uploaded")
        http.Redirect(w, r, huntAdminPath+"/?user_error=no file uploaded, please retry.", http.StatusFound)
	panic("no hunt_json uploaded")
     }
     if len(huntJSONBlobInfos) > 1 {
     	// more than one hunt_json file uploaded, delete them all (earlier defer will do this)
        http.Redirect(w, r, huntAdminPath+"/?user_error=multiple files uploaded, please only upload a single JSON file.", http.StatusFound)
	panic("multiple hunt_json files uploaded")
     }

     // check that hunt name is valid and unique
     huntname := other["hunt_name"][0]
     // only alphanumeric characters and _ are allowed
     goodNameChars, err := regexp.MatchString("^[a-zA-Z0-9_]+$", huntname)
     if !goodNameChars {
        http.Redirect(w, r, huntAdminPath+"/?user_error=Invalid hunt name. Only alphanumeric characters or underscore are allowed. Please retry.", http.StatusFound)
	panic("invalid name: "+huntname)
     }
     if err != nil {
     	serveError(c, w, err)
     	panic("error running regexp to check hunt name")
     }
     

     q := datastore.NewQuery(huntDirectoryDatastore).Filter("HuntName=", huntname).Order("-CreatedDate").Limit(1);
     q.Run(c)
     count, err := q.Count(c)	
     if err != nil {
     	// some error running query
        http.Redirect(w, r, huntAdminPath+"/?user_error=error querying hunt directory, please retry", http.StatusFound)
	panic("error querying hunt directory")
     }
     if count > 0 {
     	// hunt name already exists
        http.Redirect(w, r, huntAdminPath+"/?user_error="+"hunt name: "+huntname+" exists, please choose another.", http.StatusFound)
	panic("huntname "+huntname+" exists")
     }

     // Successful upload, check Hunt JSON 
     huntJSONBlobInfo := huntJSONBlobInfos[0]
     huntJSONReader := blobstore.NewReader(c, huntJSONBlobInfo.BlobKey)
     _, err = DecodeHuntData(huntJSONReader)
     if err != nil {
     	// Bad Hunt JSON, delete blobKey and return error
	userError := "Bad Hunt JSON: "+err.String()
	// TODO: could grab contents and pop them into JSONLint 
	err = blobstore.Delete(c, huntJSONBlobInfo.BlobKey)
	if err != nil {
	   userError += " and a further error deleting file: "+err.String()
	}
     	http.Redirect(w, r, huntAdminPath+"/?user_error="+userError, http.StatusFound)
	panic("Bad Hunt JSON uploaded")
     }
     // Good Hunt JSON, store blobKey in datastore hunt directory
    hde := HuntDirectoryEntry {
    	HuntName: huntname,
    	BlobKey:  huntJSONBlobInfo.BlobKey,
        CreatedDate:       datastore.SecondsToTime(time.Seconds()),
	Creator: user,
    }
    _, err = datastore.Put(c, datastore.NewIncompleteKey(huntDirectoryDatastore), &hde)
    if err != nil {
       panic("Error storing directory entry for hunt:"+huntname)
    }

//     http.Redirect(w, r, huntAdminDownloadPath+"/?blobKey="+string(huntJSONBlobInfo[0].BlobKey), http.StatusFound)
     http.Redirect(w, r, huntAdminPath+"/?status_message=JSON%20upload%20successful for "+huntname, http.StatusFound)
     return
}

func handleHuntAdminDownload(w http.ResponseWriter, r *http.Request) {
     c := appengine.NewContext(r)
     // get huntname for the requested blobkey
     var blobKey appengine.BlobKey
     blobKey = appengine.BlobKey(r.FormValue("blobKey"))
     c.Debugf("handleHuntAdminDownload: have blobKey %v", blobKey)
     hunts := make([]HuntDirectoryEntry, 0, 1)
     huntQuery := datastore.NewQuery(huntDirectoryDatastore).Filter("BlobKey=", blobKey).Limit(1)
    if _, err := huntQuery.GetAll(c, &hunts); err != nil {
       serveError(c, w, err)
       return
    }
    if len(hunts) != 1 {
        w.WriteHeader(http.StatusNoContent)
	return
    }
    huntName := hunts[0].HuntName
    w.Header().Set("Content-Disposition", "attachment; filename="+huntName+".json")
    blobstore.Send(w, appengine.BlobKey(blobKey))
}