var app = angular.module("Frack", ["ngRoute", "ngSanitize", "facebook"]);

app.config(function(FacebookProvider) {
  var fbID = 871188203021217;
  if (window.location.host == "localhost:8080") {
    fbID = 3518596602;
  }
  FacebookProvider.init(fbID);
});

app.config(function($routeProvider) {
  $routeProvider
    .when("/add", {
      controller: "EditController",
      templateUrl: "html/edit/editor.html"
    })
    .when("/edit/:name", {
      controller: "EditController",
      templateUrl: "html/edit/editor.html"
    })
    .when("/workout/:name", {
      controller: "ViewController",
      templateUrl: "html/view/view.html"
    })
    .when("/workout/:name/vma/:vma", {
      controller: "ViewController",
      templateUrl: "html/view/view.html"
    })
    .otherwise({
      templateUrl: "html/view/selection.html"
    });
});

//http://stackoverflow.com/a/36254259
app.directive('input', [function() {
  return {
    restrict: 'E',
    require: '?ngModel',
    link: function(scope, element, attrs, ngModel) {
      if (
        'undefined' !== typeof attrs.type && 'number' === attrs.type && ngModel
      ) {
        ngModel.$formatters.push(function(modelValue) {
          return Number(modelValue);
        });

        ngModel.$parsers.push(function(viewValue) {
          return Number(viewValue);
        });
      }
    }
  };
}]);
