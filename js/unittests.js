function setOutcome( result, body, color ) {
	var outcome = {};
	outcome.result = result;
	outcome.body = body;
	outcome.color = color;
	return outcome
}

function UnitTestCtrl($scope,$http,$location,$routeParams) {
	$scope.username = goscheduser;
	$scope.message = "Unit Tests Page";

	$scope.initValues = function() {
		$scope.ActivityListTestData = {"result":"Running.."};
		$scope.InsertActivityTestData = {"result":"Running.."};
		$scope.InsertEventTestData = {"result":"Running.."};
	}

  $scope.runtests = function() {
		$scope.initValues();
    console.log("running tests...");
		$scope.InsertActivityTest();
    console.log("...tests complete.");
	}
	$scope.InsertActivityTest = function() {
		$http({
            method : 'POST',
            url : '/gosched/v1/activity',
            data : 'title=Swim Meet&description=Gunn High School',
            headers : {
                'Content-Type' : 'application/x-www-form-urlencoded'
            }
        }).success(function (data) {
      if( data.title == "Swim Meet" ) {
				$scope.InsertActivityTestData = setOutcome("Passed", data, "#00FF00");
				$scope.InsertEventTest(data.id);
			} else {
				$scope.InsertActivityTestData = setOutcome("Failed", data, "#FF0000");
			}
    }).error(function (data) {
				$scope.InsertActivityTestData = setOutcome("Failed", data, "#FF0000");
		});
  };

	$scope.InsertEventTest = function(actid) {
		$http({
            method : 'POST',
            url : '/gosched/v1/activity/events',
            data : 'title=Private Lesson&description=Flip Turns&activity_id=' + actid,
            headers : {
                'Content-Type' : 'application/x-www-form-urlencoded'
            }
        }).success(function (data) {
      if( data.title == "Private Lesson" ) {
				$scope.InsertEventTestData = setOutcome("Passed", data, "#00FF00");
				$scope.ActivityListTest();
			} else {
				$scope.InsertEventTestData = setOutcome("Failed", data, "#FF0000");
			}
    }).error(function (data) {
				$scope.InsertEventTestData = setOutcome("Failed", data, "#FF0000");
		});
  };

	$scope.ActivityListTest = function() {
    $http.get('/gosched/v1/activity/list').success(function (data) {
      if( data instanceof Array) {
				$scope.ActivityListTestData = setOutcome("Passed", data, "#00FF00");
			} else {
				$scope.ActivityListTestData = setOutcome("Failed", data, "#FF0000");
			}
    });
  };
	
	$scope.runtests();
}

angular.module('gosched',[]).
  config(['$routeProvider', function($routeProvider) {
  $routeProvider.
      when('/tests', {templateUrl: '/partials/unittest.html', controller: UnitTestCtrl}).
			otherwise({redirectTo: '/tests'});
}]);
