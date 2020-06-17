# hugo-wrapper

This is a wrapper for the hugo command, it allows to use different version of hugo without struggle.

## Usage
### windows
example:
```bash
hugo-wrapper.exe [hugo_cmd] --hugo-version 0.72.3 [hugo_args]
``` 
hugo-version can take several form: [v]major[.minor[.patch]][-extended], or latest[-extended].
If only the major is given, the latest minor for that major will be used,
if major and minor are given, the latest patch for this major.minor will be used.
If not declared, the latest version (non extended) will be fetched.