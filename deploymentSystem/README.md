Building a vercel-like deployment system in Go

 Accepts GitHub repo  
 Spins up Docker container  
 Runs build  
 Uploads to S3 (MinIO)  
 Serves via NGINX  

Just added retry logic + deployment logs  

This is getting real.

This is exactly the kind of DIY approach that makes learning deployment systems so valuable. A few things that come to mind for your stack:

- Build verification: add checksums for uploaded artifacts
- Rollback strategy: keep last N images in MinIO with a manifest
- Health checks: NGINX should probe the container before routing