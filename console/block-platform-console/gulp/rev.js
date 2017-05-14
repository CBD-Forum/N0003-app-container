// Copyright [2016] [Lele Guo]
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
// http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// 
//**********************************
// production build script
// - depend on dev tasks
//**********************************
'use strict';

/*
var $ = {
	if: require('gulp-if'),
	notify: require('gulp-notify'),
	rev: require('gulp-rev'),
	revReplace: require('gulp-rev-replace'),
	useref: require('gulp-useref'),
	filter: require('gulp-filter'),
	uglify: require('gulp-uglify'),
	minifyCss: require('gulp-minify-css'),
	del: require('del'),
	path: require('path'),
	connect: require('connect'),
	serveStatic: require('serve-static'),
	log: require('gulp-load-plugins')()
		.loadUtils(['log'])
		.log
};
*/

// Rev:
// http://www.stevesouders.com/blog/2008/08/23/revving-filenames-dont-use-querystring/
// http://www.miwurster.com/add-gits-revision-information-to-your-app-using-gulp/
//https://segmentfault.com/q/1010000002503955/a-1020000002609153/revision

var conf = require('./conf');
var util = require('./util');
var path = require('path');
var gulp = require('gulp');
var browserSync = require('browser-sync');
var $ = require('gulp-load-plugins')({
    pattern: ['gulp-*', 'main-bower-files', 'uglify-save-license', 'del']
});
var wiredep = require('wiredep').stream;
var _ = require('lodash');


var RevAll = require('gulp-rev-all');

gulp.task('test:build:rev', ['dev'], function () {
    var revAll = new RevAll();
    return gulp.src([
        path.join(conf.paths.src, '/app/**/*.html'),
        path.join(conf.paths.src, '/assets/images/**'),
        path.join(conf.paths.tmp, '/serve/app/**')
    ]).pipe(revAll.revision())
        .pipe(gulp.dest(path.join(conf.paths.tmp, '/rev')))
        .pipe(revAll.manifestFile())
        .pipe(gulp.dest(path.join(conf.paths.tmp, '/rev')))
});

gulp.task('test:build:rev:partials', function () {
    return gulp.src([
        path.join(conf.paths.src, '/app/**/*.html'),
        path.join(conf.paths.tmp, '/serve/app/**/*.html')
    ])
        .pipe($.minifyHtml({
            empty: true,
            spare: true,
            quotes: true
        }))
        .pipe($.angularTemplatecache('templateCacheHtml.js', {
            module: conf.module,   //TODO move in one place
            root: 'app'  //TODO check
        }))
        .pipe(gulp.dest(path.join(conf.paths.tmp, '/partials')));
});


gulp.task('rev2', function () {
    var options = {
        dontRenameFile: [/^\/favicon.ico$/g, '.html', /^\/fonts\/.*/, /^\/scripts\/.*/, '.json'],
        dontSearchFile: [/^\/scripts\/.*/]
    };
    var revAll = new RevAll();

    return gulp.src([
        path.join(conf.paths.build, '/images/**'),
        path.join(conf.paths.build, '/styles/**'),
        path.join(conf.paths.build, 'index.html')
    ])
        .pipe(revAll.revision())
        .pipe(gulp.dest('build/final'))
        .pipe(revAll.manifestFile())
        .pipe(gulp.dest('build/final'));
});

gulp.task('rev', function () { // ['build'],

    //dontSearchFile : [ /^\/scripts\/.*/ ]
    var options = { dontRenameFile: [/^\/favicon.ico$/g, '.html', /^\/fonts\/.*/, '.json'] };

    options.annotator = function (contents, path) {
        var fragments = [{ 'contents': contents, 'path': path }];
        return fragments;
    };

    options.replacer = function (fragment, replaceRegExp, newReference, referencedFile) {
        if (fragment.path === referencedFile.revPathOriginal) {
            return;
        }
        fragment.contents = fragment.contents.replace(replaceRegExp, '$1' + newReference + '$3$4');
    };

    var revAll = new RevAll(options);

    return gulp.src([path.join(conf.paths.build, '/**')])
        .pipe(revAll.revision())
        .pipe(gulp.dest('build/final'))
        .pipe(revAll.manifestFile())
        .pipe(gulp.dest('build/final'));
});

gulp.task('rv', function () {

    var f = $.filter(['!build/fonts', '!build/locales'], { restore: true });

    return gulp.src([path.join(conf.paths.build, '/**')
    ])
        .pipe(f)
        .pipe($.rev())
        .pipe(gulp.dest(path.join(conf.paths.build, '/rev')))
        .pipe($.rev.manifest())
        .pipe(gulp.dest(path.join(conf.paths.build, '/rev')))
})

gulp.task('rvv', ['rv'], function () {
    var manifest = gulp.src(path.join(conf.paths.build, '/rev', '/rev-manifest.json'));

    return gulp.src(path.join(conf.paths.src, 'index.html'))
        .pipe($.revReplace({ manifest: manifest }))
        .pipe(gulp.dest(path.join(conf.paths.build, '/revv')));
});




gulp.task('run:rev', function () {
  //browserSyncInit(path.join(conf.paths.build,'/final'));
});





gulp.task('build:rename', function () {
    return new Promise(function (resolve, reject) {
        var vp = vinylPaths();

        gulp.src([path.join(conf.paths.build, '**/*.js'), path.join(conf.paths.build, '**/*.css')])
            .pipe(vp)
            .pipe($.rename({
                suffix: '_v' + conf.pkg.version
            }))
            .pipe(gulp.dest(path.join(conf.paths.build, '/')))
            .on('end', function () {
                $.del(vp.paths).then(resolve).catch(reject);
            });
    });
   
    // return gulp.src([path.join(conf.paths.build, 'scripts/**/*.js'), path.join(conf.paths.build, 'styles/**/*.css')])
    //     .pipe($.rename({
    //         suffix: '_v' + conf.pkg.version
    //     }))
    //     .pipe(vinylPaths(function (paths) {
    //         console.log('Paths:', paths);
    //         return Promise.resolve();
    //     }))
    
    // .pipe($.rename(function (filepath) {
    //     var origin= filepath.basename+filepath.extname;
    //     if (filepath.extname == '.js') {
    //         filepath.dirname += "/scripts";
    //     } else {
    //         filepath.dirname += "/styles";
    //     }
    //     var resourceDir = path.join(__dirname, '..', conf.paths.build, filepath.dirname);
    //     console.log(path.join(resourceDir,origin))
           
    //     filepath.basename += '_v' + conf.pkg.version;
    //     // fs.unlinkSync(path.join(resourceDir,origin))
    //     //$.del(path.join(resourceDir,origin));
    //     return filepath; 
    // }))
    //    .pipe(gulp.dest(path.join(conf.paths.build, '/')));
});




//http://stackoverflow.com/questions/23820703/how-to-inject-content-of-css-file-into-html-in-gulp
gulp.task('xx:build:rev', function () {

    return gulp.src(path.join(conf.paths.build, 'index.html'))
        .pipe($.replace(/(styles|scripts)\/([^\.]+\.(css|js))/g, function (match, dir, filename, extname) {
            // rename file BUT also need rename the real file in styles & script dir
            //dir: styles|scripts
            //filename:  xxx.js | xxx.css
            //extname:  css |  js
            var resourceDir = path.join(__dirname, '..', conf.paths.build, dir);
            var filebasename = path.basename(filename, '.' + extname);
            var revFilename = filebasename + '_v' + ver + '.' + extname;
            fs.rename(path.join(resourceDir, filename), path.join(resourceDir, revFilename), function (err) {
                //gutil.log(gutil.colors.red('[Rename]'), err);
            });
            console.log(path.join(dir, revFilename));
            return path.join(dir, revFilename);
        })).on('error', conf.errorHandler('Revision'))
    /* 
          .pipe($.replace(/(<link rel="stylesheet" href="styles\/)([^\.]+)(\.css"[^>]*>)/g, function (match, p1, p2, p3) {
              var stylesDir = path.join(__dirname, '..', conf.paths.build, 'styles')
              var revisoned = p2 + '_v' + ver;
              var styleFile = path.join(stylesDir, p2 + '.css');
              var styleFileRev = path.join(stylesDir, revisoned + '.css');
              //fs.renameSync
              fs.rename(styleFile, styleFileRev, function (err) {
                  gutil.log(gutil.colors.red('[Rename]'), err.toString());
              });
              return p1 + revisoned + p3;
          }))
          */
    // .pipe($.replace(/(<link rel="stylesheet" href="styles\/)([^\.]+)(\.css"[^>]*>)/g, '$1$2_' + ver + '$3'))
        .pipe(gulp.dest(path.join(conf.paths.build, '/')));
});
