# Versions

### A tool to generate a simple version-lapse (timelapse) for a front-end in a repo (git repository)

## Why

I always think it's fun to look at the progess of any project that I - or others - make. Sadly, it's a bore to check out every version committed on git and starting a server. Thus, this tool has come to life! 🤓

## What

Explained in simple words: It makes a version-lapse of a repo.

Please checkout the gif below to see my example of this.

![https://voters.cafe version-lapse of five commits](sample/out.gif?raw=true "https://voters.cafe version-lapse of five commits")

## How

It's quite simple to use. I promise! Download the source-code of this repo. Fire up a terminal and type in `go run . -repo ../path/to/git/repo -commits 5 -wait 5 -port 5000`, this will make a version-lapse of of the last `5 commits` of the repo `../path/to/git/repo`. It starts a dev-server for each commit - from ports `5000` 'til `5004`, waits `5 seconds`, takes a screendump, closes the dev-server and annotates the image. Lastly it makes the `5 commits` into a gif. Easy peasy!

## Usage

First, make sure you have `ImageMagick <= 7.x` installed.

Second, have a `Chromium`-browser installed.

Third, and most importantly, have `go >= 1.11` installed - As I'm not yet providing native builds of this program.

And, of course, `git` and a local repo.

Then we have some flags to configure.

Required:

 * `-repo` - specifies the path to the repo you'd like to make a verison-lapse of.

Useful:

 * `-dump`
   * defaults to `./screendumps`, is the directory to dump the screenshots in.

 * `-out`
   * defaults to `./out.gif`, is the path and name to place the generated gif.

 * `-commits`
   * defaults to `0`, meaning all the commits. 

 * `-port`
   * defaults to `5000`, and will use ports from `5000` to `5000 + N_COMMITS`.

 * `-wait`
   * defaults to `5`, this will let the program wait for `5` seconds to start the app, `5` seconds to load the page and (always) `1` second to process the screendump.
   * This adds up to `5 + 5 + 1 = 11` seconds for each commit to process.
   * Or `55` seconds - about a minute - to run through `5` commits.

## Post Scriptum

Only works with `yarn` and a `dev`-script in your `package.json`.

If the program crashes at any point, you'd `git checkout master` on the repo specified. The program **will** try to clean up after itself, but might fail at times.

If you by any reason decide to run the program again on the same repo or any other repo, it'll clean the `screendumps`-dir and `out.gif`-file. So, make sure to save whatever you need!

*Thanks!* - Mads Cordes ([@Mobilpadde](https://twitter.com/Mobilpadde "Twitter")).