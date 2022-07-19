// https://laravel-mix.com/docs/main/what-is-mix

let fs = require("fs");
let glob = require("glob");
let path = require("path");
let exec = require("child_process").execSync;
let mix = require("laravel-mix");
let { optimize } = require('svgo');

let ImageminPlugin = require('imagemin-webpack-plugin').default;

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


mix.extend("svg", function (_webpackConfig, src, dest) {
    mix.copy(src, dest);

    // optimize svg files
    glob(dest + "/**/*.svg", (err, files) => {
        files.forEach(file => {
            let svg = fs.readFileSync(file, "utf8");
            let result = optimize(svg);
            fs.writeFileSync(file, result.data);
        });
    });
})

// Either copies the templates from the given source to the given destination
// (prod) or creates a symlink at the destination to the source (dev).
mix.extend("copyTemplates", function (_webpackConfig, src, dest) {
    if (mix.inProduction()) {
        // in prod mode we copy the templates to the output path
        mkdirs(path.dirname(dest))
        try {
            // check if path already exists
            let stat = fs.lstatSync(dest)
            if (stat.isSymbolicLink()) {
                fs.unlinkSync(dest)
            }
        } catch (err) {
            if (err.code && err.code !== "ENOENT") {
                throw err;
            }
        }
        mix.copy(src, dest);

    } else {
        // in dev mode we create a symlink to the template
        mkdirs(path.dirname(dest))

        src = "../../" + src  // TODO: maybe find a better way to do this
        try {
            // check if path already exists
            let stat = fs.lstatSync(dest)
            if (stat.isSymbolicLink()) {
                // check if symlink is correct
                let target = fs.readlinkSync(dest)
                if (target === src) {
                    return
                }
                console.log(`Removing invalid symlink ${dest}`)
                fs.unlinkSync(dest)
            } else if (stat.isDirectory()) {
                console.log(`Removing directory ${dest}`)
                fs.rmdirSync(dest, { recursive: true });
            }
        } catch (err) {
            if (err.code && err.code !== "ENOENT") {
                throw err;
            }
        }

        fs.symlinkSync(src, dest)
    }
})

// Generates the assets.go file inside the assets directory that is responsible
// for loading/embedding the assets in the app.
mix.extend("genAssetsGo", function (_webpackConfig, dest) {
    if (mix.inProduction()) {
        exec(`go run resources/gen_assets_fs.go prod '${dest}'`)
    } else {
        exec(`go run resources/gen_assets_fs.go dev '${dest}'`)
    }
})

class Imagemin {
    dependencies() {
        this.requiresReload = `
            Imagemin's required plugins have been installed.
            Please run "npm run dev" again.
        `;
        return ['copy-webpack-plugin', 'imagemin-webpack-plugin'];
    }

    register(context, from, to, imageminOptions = {}) {
        imageminOptions.externalImages === undefined;
        let externalImages = imageminOptions.externalImages || {};
        externalImages.sources = glob.sync(from);
        externalImages.context = context;  // context is the path prefix of sources that will be stripped off
        externalImages.destination = to;
        imageminOptions.externalImages = externalImages;

        this.tasks = this.tasks || [];
        this.tasks.push(new ImageminPlugin(imageminOptions));
    }

    webpackPlugins() {
        return this.tasks;
    }
}

mix.extend('imagemin', new Imagemin());

mix.sourceMaps()
    // fonts & icons
    .copy(glob.sync("node_modules/@fontsource/roboto/files/roboto-latin-*"), staticPath + "/fonts/roboto")
    .imagemin( // copy and minify icons
        "node_modules/@tabler/icons/icons",
        "node_modules/@tabler/icons/icons/*.svg",
        outputPath + "/icons",
        {
            cacheFolder: outputPath + "/.cache",
            svgo: {
                plugins: [
                    { removeViewBox: false },
                ]
            },
        }
    )
    .sass("resources/scss/tabler-icons.scss", "css")

    // tabler.io
    .js("resources/js/tabler.js", "js")
    .sass("resources/scss/tabler.scss", "css")
    .sass("resources/scss/tabler-vendors.scss", "css")

    // libs

    // app
    .copyTemplates("resources/views", outputPath + "/templates")
    .imagemin( // copy and minify images
        "resources",
        "resources/img/**/*(*.svg|*.jpg|*.png)",
        staticPath,
        {
            cacheFolder: outputPath + "/.cache",
            svgo: {
                plugins: [
                    { removeViewBox: false },
                ]
            },
            optipng: {
                optimizationLevel: 5
            },
        }
    )
    .sass("resources/scss/app.scss", "css")
    .js("resources/js/app.js", "js")
    .js("resources/js/demo.js", "js")
    .sass("resources/scss/demo.scss", "css")
    .genAssetsGo(outputPath + "/assets.go");
