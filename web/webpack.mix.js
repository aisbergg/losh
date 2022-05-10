// https://laravel-mix.com/docs/main/what-is-mix

let mix = require("laravel-mix");
let glob = require("glob");
let fs = require("fs");
let path = require("path");

let outputPath = "build/assets"
let staticPath = path.join(outputPath, "static")

function mkdirs(dir) {
    comps = dir.split(path.sep)
    if (comps.length === 0) {
        return
    }
    cur = ""
    for (var i = 0; i < comps.length; i++) {
        cur = path.join(cur, comps[i])
        if (fs.existsSync(cur)) {
            if (!fs.statSync(cur).isDirectory()) {
                throw new Error(`Not a directory: "${cur}"`);
            }
        } else {
            fs.mkdirSync(cur);
        }
    }
}

mix.options({
    fileLoaderDirs: {
        images: "img",
        fonts: "fonts"
    },
    processCssUrls: false,
    publicPath: staticPath,  // base output path for generated assets
    resourceRoot: "static",  // base path used for serving the assets (used with processCssUrls)
})

mix.extend("symlink", function (_webpackConfig, src, dest) {
    if (!mix.inProduction()) {
        mkdirs(path.dirname(dest))

        try {
            // check if path already exists
            let stat = fs.lstatSync(dest)
            if (stat.isSymbolicLink()) {
                // check if symlink is correct
                let target = fs.readlinkSync(dest)
                if (target === src) {
                    return
                }
                fs.unlinkSync(dest)
            }
        } catch (err) {
            if (err.code && err.code !== "ENOENT") {
                throw err;
            }
        }

        fs.symlinkSync(src, dest)
    }
})

mix.sourceMaps()
    // fonts & icons
    .copy(glob.sync("node_modules/@fontsource/roboto/files/roboto-latin-*"), staticPath + "/fonts/roboto")
    .copy("node_modules/@tabler/icons/iconfont/fonts", staticPath + "/fonts/iconfont")
    .sass("resources/scss/tabler-icons.scss", "css")

    // tabler.io
    .js("resources/js/tabler.js", "js")
    .sass("resources/scss/tabler.scss", "css")
    .sass("resources/scss/tabler-vendors.scss", "css")

    // libs

    // app
    .symlink("../../resources/views", outputPath + "/templates")
    .copyDirectory("resources/img", staticPath + "/img")
    .sass("resources/scss/app.scss", "css")
    .js("resources/js/app.js", "js")
    .js("resources/js/demo.js", "js")
    .sass("resources/scss/demo.scss", "css");
