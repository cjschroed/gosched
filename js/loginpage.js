
function LoginPageCtrl($scope,$http,$location,$routeParams) {
	console.log("Login Page controller");
}

angular.module('gosched', []).
  config(['$routeProvider', function($routeProvider) {
  $routeProvider.
      when('/', {templateUrl: '/partials/loginpage.html', controller: LoginPageCtrl}).
			otherwise({redirectTo: '/'});
}]);
