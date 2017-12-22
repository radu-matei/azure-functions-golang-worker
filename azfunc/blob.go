package azfunc

import (
	"time"
)

// the following structs are taken from the Azure Storage Go library - https://github.com/Azure/azure-storage-go/blob/master/blob.go
// not importing the package yet as it deals with xml - added json tags to automatically unmarshal

// BlobProperties contains various properties of a blob
// returned in various endpoints like ListBlobs or GetBlobProperties.
type BlobProperties struct {
	LastModified       time.Time `xml:"Last-Modified" json:"LastModified"`
	Etag               string    `xml:"Etag" json:"ETag"`
	ContentMD5         string    `xml:"Content-MD5" header:"x-ms-blob-content-md5" json:"ContentMD5"`
	ContentLength      int64     `xml:"Content-Length" json:"Length"`
	ContentType        string    `xml:"Content-Type" header:"x-ms-blob-content-type" json:"ContentType"`
	ContentEncoding    string    `xml:"Content-Encoding" header:"x-ms-blob-content-encoding" json:"ContentEncoding"`
	CacheControl       string    `xml:"Cache-Control" header:"x-ms-blob-cache-control" json:"CacheControl"`
	ContentLanguage    string    `xml:"Cache-Language" header:"x-ms-blob-content-language" json:"ContentLanguage"`
	ContentDisposition string    `xml:"Content-Disposition" header:"x-ms-blob-content-disposition" json:"ContentDisposition"`
	BlobType           BlobType  `xml:"x-ms-blob-blob-type" json:"BlobType"`
	SequenceNumber     int64     `xml:"x-ms-blob-sequence-number" json:"PageBlobSequenceNumber"`
	// TODO these properties do not appear in the json or in the CloudBlob type from .NET - investigate if needed or can be removed
	CopyID                string    `xml:"CopyId"`
	CopyStatus            string    `xml:"CopyStatus"`
	CopySource            string    `xml:"CopySource"`
	CopyProgress          string    `xml:"CopyProgress"`
	CopyCompletionTime    time.Time `xml:"CopyCompletionTime"`
	CopyStatusDescription string    `xml:"CopyStatusDescription"`
	//
	LeaseStatus     string `xml:"LeaseStatus" json:"LeaseStatus"`
	LeaseState      string `xml:"LeaseState" json:"LeaseState"`
	LeaseDuration   string `xml:"LeaseDuration" json:"LeaseDuration"`
	ServerEncrypted bool   `xml:"ServerEncrypted" json:"IsServerEncrypted"`
	IncrementalCopy bool   `xml:"IncrementalCopy" json:"IsIncrementalCopy"`
}

// BlobType defines the type of the Azure Blob.
type BlobType string

// Types of page blobs
const (
	BlobTypeBlock  BlobType = "BlockBlob"
	BlobTypePage   BlobType = "PageBlob"
	BlobTypeAppend BlobType = "AppendBlob"
)

// Sys - to clarify the name / use of this type
type Sys struct {
	MethodName string    `json:"MethodName"`
	UtcNow     time.Time `json:"UtcNow"`
	RandGUID   string    `json:"RandGuid"`
}

// Blob contains the data from a blob as string
type Blob struct {
	Data string
}
