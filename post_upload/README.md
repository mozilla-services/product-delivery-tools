# Post Upload
Copies a list of files to their S3 release locations.

```
NAME:
   post_upload - 

USAGE:
   post_upload [global options] command [command options] [arguments...]

VERSION:
   1.0.0

AUTHOR(S): 
   Jeremy Orem <oremj@mozilla.com> 

COMMANDS:
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --product, -p 			Set product name to build paths properly.
   --version, -v 			Set version number to build paths properly.
   --nightly-dir "nightly"		Set the base directory for nightlies (ie $product/$nightly_dir/}, and the parent directory for release candidates (default 'nightly'}.
   --branch, -b 			Set branch name to build paths properly.
   --buildid, -i 			Set buildid to build paths properly.
   --build-number, -n 			Set buildid to build paths properly.
   --revision, -r 			
   --who, -w 				
   --no-shortdir, -S			Don't symlink the short dated directories.
   --builddir 				Subdir to arrange packaged unittest build paths properly.
   --subdir 				Subdir to arrange packaged unittest build paths properly.
   --tinderbox-builds-dir 		Set tinderbox builds dir to build paths properly.
   --release-to-latest, -l		Copy files to $product/$nightly_dir/latest-$branch
   --release-to-dated, -d		Copy files to $product/$nightly_dir/$datedir-$branch
   --release-to-candidates-dir, -c	Copy files to $product/$nightly_dir/$version-candidates/build$build_number
   --release-to-mobile-candidates-dir	Copy mobile files to $product/$nightly_dir/$version-candidates/build$build_number/$platform
   --release-to-tinderbox-builds, -t	Copy files to $product/tinderbox-builds/$tinderbox_builds_dir
   --release-to-latest-tinderbox-builds	Softlink tinderbox_builds_dir to latest
   --release-to-tinderbox-dated-builds	Copy files to $product/tinderbox-builds/$tinderbox_builds_dir/$timestamp
   --release-to-try-builds		Copy files to try-builds/$who-$revision
   --signed				Don't use unsigned directory for uploaded files
   --help, -h				show help
```
