# geosite_editor
simple tool for editing geosite.dat

# Run
```bash
./geosite_editor -t tasks.json
```

# Example
```json
{
  "src": "geosite.dat",
  "dst": "mygeosite.dat",
  "tasks": [
    {
      "type": "remove",
      "country_code": "US",
      "domains": [
        "tiktok.com"
      ]
    },
    {
      "type": "add",
      "country_code": "GFW",
      "domains": [
        "csdn.net",
        "zhihu.com"
      ]
    },
    {
      "type": "remove",
      "country_code": "*",
      "domains": [
        "cloudflare.com"
      ]
    },
    {
      "type": "copy",
      "src_country_code": "GFW",
      "dst_country_code": "AD"
    },
    {
      "type": "delete",
      "entries": [
        "CA",
        "CLOUDFLARE"
      ]
    }
  ]
}
```
