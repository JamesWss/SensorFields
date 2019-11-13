var gulp = require('gulp');
var less = require('gulp-sass');
var path = require('path');
var shell = require('gulp-shell');

var lessPath = './src/distributed/webclient/sass/**/*.sass';
var goPath = 'src/distributed/**/*.go';

gulp.task('less', function () {
  return gulp.src('./src/distributed/webclient/sass/app.sass')
    .pipe(less({
      paths: [ path.join(__dirname, 'sass', 'includes') ]
    }))
    .pipe(gulp.dest('./src/distributed/web/res/css'));
});

gulp.task('compilepkg', function() {
  return gulp.src(goPath, {read: false})
    .pipe(shell(['go install <%= stripPath(file.path) %>'],
      {
          templateData: {
            stripPath: function(filePath) {
              var subPath = filePath.substring(process.cwd().length + 5);
              var pkg = subPath.substring(0, subPath.lastIndexOf(path.sep));
              return pkg;
            }
          }
      })
    );
});

gulp.task('watch', function() {
  gulp.watch(lessPath, ['less']);
  gulp.watch(goPath, ['compilepkg']);
});
