# GOWORK

## Features

### Image compression
Compress images before save in AWS S3

1. Client send the file path and bucket and the new path and bucket
2. Server compress the image
3. Server save the compressed image in AWS S3 bucket

#### Compression options
```js
{
  "bucket": "my-bucket",
  "path": "/path/to/image.jpg",
  "newBucket": "my-new-bucket", // if not provided, the original bucket will be used (optional)
  "newPath": "/path/to/new/image.jpg", // if not provided, the original path will be used (optional)
  "quality": 80 // (0-100) (default: 80) (optional)
  "contentType": "image/jpeg", // (default: image/jpeg) (optional)
  "deleteOriginal": true // (default: false) (optional)
}
```
---

### Send email
Send an email to a user

1. Client send the email destination, subject and template id
2. Server get the template from AWS S3 and send the email

```js
{
  "to": "email@example.com",
  "subject": "Hello",
  "templatePath": "path/to/template.html",
  "data": {
    "name": "John Doe"
  }
}
```