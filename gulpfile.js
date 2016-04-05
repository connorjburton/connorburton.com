var gulp = require('gulp');
var sass = require('gulp-sass');
var autoprefixer = require('gulp-autoprefixer');

gulp.task('css', function() {
	return gulp.src('./css/**/*.scss')
		.pipe(autoprefixer())
		.pipe(sass().on('error', sass.logError))
		.pipe(gulp.dest('./public/css'));
});

gulp.task('watch', function() {
	gulp.watch('./css/**/*.scss', ['css']);
});