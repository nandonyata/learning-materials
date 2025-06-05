# publish-container-image-to-ghcr
Source: https://www.youtube.com/watch?v=RgZyX-e6W9E

1. Login to GHCR using Github Personal Access Token (classic) (PAT)
```
docker login --username <GITHUB_USERNAME> --password <GITHUB_PAT> ghcr.io
```

2. Build this code as an image
```
docker build . -t ghcr.io/<GITHUB_USERNAME>/<PACKAGE_NAME>:latest
```

3. Push
```
docker push ghcr.io/<GITHUB_USERNAME>/<PACKAGE_NAME>:latest
```

4. Run and test the image
```
docker run ghcr.io/<GITHUB_USERNAME>/<PACKAGE_NAME>:latest
```

Example: 
```
docker run ghcr.io/nando10xers/hello-world-ghcr:latest
```