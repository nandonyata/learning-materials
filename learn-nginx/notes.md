Source: https://www.youtube.com/watch?v=q8OleYuqntY

Nginx uses for:
- Web server
- Load balancer
- Caching
- Encrypted communication (security)
- Compression response
- segmentation

In production, reccommended to set nginx worker_processes value to be equal to cpu cores of the server
where nginx is running, or set to "auto".