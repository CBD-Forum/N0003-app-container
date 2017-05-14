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
// devlopment build script
//**********************************
'use strict';

var conf = require('./conf');
var path = require('path');
var fs = require('fs');
var gulp = require('gulp');
var gutil = require('gulp-util');
var browserSync = require('browser-sync');
var reload = browserSync.reload;
var $ = require('gulp-load-plugins')({
    pattern: ['gulp-*', 'del']
});

var runSequence = require('run-sequence');
var wiredep = require('wiredep').stream;
var htmlReporter = require('gulp-csslint-report');
var esreporter = require('eslint-html-reporter');
var _ = require('lodash');

/**
 * Compile vendor less to css to tmp serve dir.
 *
 * NOTE: this is only used for know how to generate the vendor css.[ NOT USED ]
 */
gulp.task('vendor:less', function () {
    var lessOptions = {
        paths: [
            conf.wiredepOptions.directory
        ]
    };

    return gulp.src([
        path.join(conf.paths.src, '/app/vendor.less')
    ])
        .pipe(wiredep(conf.wiredepOptions))
        .pipe($.sourcemaps.init())
        .pipe($.less(lessOptions)).on('error', conf.errorHandler('Less'))
        .pipe($.uncss({
            html: [
                path.join(conf.paths.src, '/app/**/*.html')
            ]
        }))
        .pipe($.autoprefixer()).on('error', conf.errorHandler('Autoprefixer'))
        .pipe($.sourcemaps.write())
        .pipe(gulp.dest(path.join(conf.paths.tmp, '/serve/app/')));
    /*
	  return gulp
		.src(config.paths.vendorcss)
		.pipe(concat('vendor.min.css'))
		.pipe(bytediff.start())
		.pipe(minifyCss())
		.pipe(bytediff.stop(bytediffFormatter))
		.pipe(gulp.dest(config.paths.destination));
        */
});

//*****************************************************************
// LESS
//*****************************************************************


/**
 * TODO create rules
 *
 * https://github.com/CSSLint/csslint/wiki/Rules-by-ID
 */
gulp.task('app:csslint', function () {
    return gulp.src([
        path.join(conf.paths.src, '/app/**/*.less'),
        path.join(conf.paths.src, '/assets/styles/less/app.less'),
        path.join('!' + conf.paths.src, '/app/app.less')
    ], { read: false })
        .pipe($.csslint())
        .pipe($.csslint.reporter());
    //.pipe($.csslint.reporter('fail')); // Fail on error (or csslint.failReporter())
});

/**
 * Compile app less to css to tmp serve dir
 */
gulp.task('app:less', function () {
    //less
    var lessOptions = {
        paths: [
            path.join(conf.paths.src, '/app')
        ]
    };

    var injectFiles = gulp.src([
        path.join(conf.paths.src, '/app/**/*.less'),
        path.join(conf.paths.src, '/assets/styles/less/app.less'),
        path.join('!' + conf.paths.src, '/app/app.less')
    ], { read: false });

    var injectOptions = {
        transform: function (filePath) {
            filePath = filePath.replace(conf.paths.src + '/app/', '');
            return '@import "' + filePath + '";';
        },
        starttag: '// injector',
        endtag: '// endinjector',
        addRootSlash: false
    };

    return gulp.src([
        path.join(conf.paths.src, '/app/app.less')
    ])
    //less lint
    //.pipe($.recess())
    //.pipe($.recess.reporter()).on('error', conf.errorHandler('recess'))
        .pipe($.inject(injectFiles, injectOptions))
        .pipe($.sourcemaps.init())
        .pipe($.less(lessOptions)).on('error', conf.errorHandler('Less'))
     
    // *********************************
    //.pipe($.uncss({
    //    html: [
    //       path.join(conf.paths.src, '/app/**/*.html')
    //    ]
    // }))
    // *********************************
    //css lint
    //.pipe($.csslint())
    //.pipe(htmlReporter({
    //     filename: 'csslint-report.html',
    //     directory: path.join('.', conf.paths.lint, '/')
    // }))
    //.pipe($.csslint.reporter())
        .pipe($.autoprefixer()).on('error', conf.errorHandler('Autoprefixer'))
        .pipe($.sourcemaps.write())
        .pipe(gulp.dest(path.join(conf.paths.tmp, '/serve/app/')))
        .pipe(browserSync.stream());
    //.pipe(reload({stream:true})); //OR this if the above dosen't work
});

gulp.task('app:scss', function () {
    //sass
    var sassOptions = {
        style: 'expanded'
    };

    var injectFiles = gulp.src([
        path.join(conf.paths.src, '/app/**/*.scss'),
        path.join(conf.paths.src, '/assets/styles/scss/app.scss'),
        path.join('!' + conf.paths.src, '/app/app.scss')
    ], { read: false });

    var injectOptions = {
        transform: function (filePath) {
            filePath = filePath.replace(conf.paths.src + '/app/', '');
            return '@import "' + filePath + '";';
        },
        starttag: '// injector',
        endtag: '// endinjector',
        addRootSlash: false
    };

    return gulp.src([
        path.join(conf.paths.src, '/app/app.scss')
    ])
    //less lint
    //.pipe($.recess())
    //.pipe($.recess.reporter()).on('error', conf.errorHandler('recess'))
        .pipe($.inject(injectFiles, injectOptions))
        .pipe($.sourcemaps.init())
        .pipe($.sass(sassOptions)).on('error', conf.errorHandler('Sass'))
    // *********************************
    //.pipe($.uncss({
    //    html: [
    //       path.join(conf.paths.src, '/app/**/*.html')
    //    ]
    // }))
    // *********************************
    //css lint
    //.pipe($.csslint())
    //.pipe(htmlReporter({
    //     filename: 'csslint-report.html',
    //     directory: path.join('.', conf.paths.lint, '/')
    // }))
    //.pipe($.csslint.reporter())
        .pipe($.autoprefixer({
            browsers: ['last 2 versions'],
            cascade: false
        })).on('error', conf.errorHandler('Autoprefixer'))
        .pipe($.sourcemaps.write())
        .pipe(gulp.dest(path.join(conf.paths.tmp, '/serve/app/')))
        .pipe(browserSync.stream());
    //.pipe(reload({stream:true})); //OR this if the above dosen't work
});


//*****************************************************************
// Javascript
//*****************************************************************


// JSHint
function jshint() {
    return gulp.src(path.join(conf.paths.src, '/app/**/*.js'))
        .pipe($.jshint())
        .pipe($.jshint.reporter('gulp-jshint-html-reporter', {
            filename: path.join(conf.paths.lint, '/jshint-report.html'),
            createMissingFolders: true
        }));
}

// ESHint
function eslint() {
    return gulp.src(path.join(conf.paths.src, '/app/**/*.js'))
    // eslint() attaches the lint output to the "eslint" property
    // of the file object so it can be used by other modules.
        .pipe($.eslint())
    // eslint.format() outputs the lint results to the console.
    // Alternatively use eslint.formatEach() (see Docs).
        .pipe($.eslint.format(esreporter, function (results) {
            var dir = path.join(__dirname, '..', conf.paths.lint)
            if (!fs.existsSync(dir)) {
                //fs.rmdirRecursive(dir);
                fs.mkdirSync(dir);
            }
            fs.writeFileSync(path.join(dir, '/eslint-report.html'), results);
        }))
    // To have the process exit with an error code (1) on
    // lint error, return the stream and pipe to failAfterError last.
        .pipe($.eslint.failAfterError());
}

var jsTaskStream = jshint;

/**
 * Process javascript, mainly process linting.
 */
gulp.task('app:js', jsTaskStream);

/**
 * Inject js and css to index.html
 */
gulp.task('inject:jscss', ['app:less', 'app:js'], function () {
    var injectStyles = gulp.src([
        path.join(conf.paths.tmp, '/serve/app/**/*.css'),
        path.join('!' + conf.paths.tmp, '/serve/app/vendor.css')
    ], { read: false });

    //Angular needs the module definitions to be registered before they are used.
    //Always grab module files first
    var injectScripts = gulp.src([
        path.join(conf.paths.src, '/app/**/*.module.js'),
        path.join(conf.paths.src, '/app/**/*.js'),
        path.join('!' + conf.paths.src, '/app/**/*.spec.js'),
        path.join('!' + conf.paths.src, '/app/**/*.mock.js'),
    ])
        .pipe($.angularFilesort()).on('error', conf.errorHandler('AngularFilesort'));

    var injectOptions = {
        //The "unixified" path to the file with any ignorePath's removed,
        ignorePath: [conf.paths.src, path.join(conf.paths.tmp, '/serve')],

        //The root slash is automatically added at the beginning of the path ('/')
        //or removed if set to false.
        addRootSlash: false
    };

    return gulp.src(path.join(conf.paths.src, '/*.html'))
        .pipe($.inject(injectStyles, injectOptions))
        .pipe($.inject(injectScripts, injectOptions))
        .pipe(wiredep(conf.wiredepOptions))
        .pipe(gulp.dest(path.join(conf.paths.tmp, '/serve')));
});


/**
 * Build dev (in sequence)
 */
gulp.task('dev', ['inject:jscss']);

/**
 * Build dev (in parallel)
 */
// THIS SHOULD BE DEPRECATED WHEN gulp 4.0 is out
// it will provide out of box: gulp.parallel(...tasks)

// This will run in this order:
// * build-scripts and build-styles in parallel
// * Finally call the callback function
gulp.task('dev:async', function (callback) {
    runSequence(['app:less', 'app:js'], function () {
        console.log('DONE');
        callback();
    });
});
