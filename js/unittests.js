function setOutcome( result, body, color ) {
	var outcome = {};
	outcome.result = result;
	outcome.body = body;
	outcome.color = color;
	return outcome
}

function UnitTestCtrl($scope,$http,$location,$routeParams) {
	$scope.username = goscheduser;
	$scope.readme = "/README.md";
	$scope.message = "Welcome to the Go Scheduling API";

	$scope.initValues = function() {
		$scope.addbook = false;
		$scope.ActivityListTestData = {"result":"Ready."};
		$scope.InsertActivityTestData = {"result":"Ready."};
		$scope.InsertEventTestData = {"result":"Ready."};
		$scope.InsertBookingTestData = {"result":"Ready."};
		$scope.DeleteEventTestData = {"result":"Ready."};
		$scope.DeleteActivityTestData = {"result":"Ready."};
		$scope.GetActivityTestData = {"result":"Ready."};
		$scope.GetEventTestData = {"result":"Ready."};
		$scope.EventSearch1TestData = {"result":"Ready."};
		$scope.EventSearch2TestData = {"result":"Ready."};
	}

  $scope.runtests = function() {
		$scope.initValues();
    console.log("running tests...");
		$scope.InsertActivityTest();
		$scope.ActivityListTest();
		$scope.DeleteEventTest();
		$scope.DeleteActivityTest();
    console.log("...tests complete.");
		$scope.addbook = true;
	}
	$scope.InsertActivityTest = function() {
		$http({
            method : 'POST',
            url : '/gosched/v1/activity',
            data : 'title=Swim Meet&description=Gunn High School&vendor_name=LAMVAC',
            headers : {
                'Content-Type' : 'application/x-www-form-urlencoded'
            }
        }).success(function (data) {
      if( data.title == "Swim Meet" ) {
				$scope.actid = data.id;
				$scope.InsertActivityTestData = setOutcome("Passed", data, "#00FF00");
				$scope.InsertEventTest();
			} else {
				$scope.InsertActivityTestData = setOutcome("Failed", data, "#FF0000");
			}
    }).error(function (data) {
				$scope.InsertActivityTestData = setOutcome("Failed", data, "#FF0000");
		});
  };

	$scope.InsertEventTest = function() {
		$http({
            method : 'POST',
            url : '/gosched/v1/activity/events',
            data : 'title=Private Lesson&description=Flip Turns&activity_id=' + $scope.actid + '&start_time=Jan 2, 2013 at 3:00pm (MST)&end_time=Jan 2, 2013 at 4:00pm (MST)&max_attendees=8',
            headers : {
                'Content-Type' : 'application/x-www-form-urlencoded'
            }
        }).success(function (data) {
      if( data.title == "Private Lesson" ) {
				$scope.InsertEventTestData = setOutcome("Passed", data, "#00FF00");
				$scope.eventid = data.id;
				$scope.InsertBookingTest();
				$scope.GetActivityTest();
				$scope.GetEventTest();
				$scope.EventSearch1Test();
				$scope.EventSearch2Test();
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

  $scope.InsertBookingTest = function() {
		console.log('adding booking for eventid: ' + $scope.eventid);
    $http({
            method : 'POST',
            url : '/gosched/v1/activity/book',
            data : 'event_id=' + $scope.eventid,
            headers : {
                'Content-Type' : 'application/x-www-form-urlencoded'
            }
      }).success(function (data) {
       	$scope.InsertBookingTestData = setOutcome("Passed", data, "#00FF00");
    	}).error(function (data) {
        $scope.InsertBookingTestData = setOutcome("Failed", data, "#FF0000");
    });
  };

	$scope.DeleteActivityTest = function() {
		$http({
            method : 'POST',
            url : '/gosched/v1/activity',
            data : 'title=Swim Meet&description=Gunn High School',
            headers : {
                'Content-Type' : 'application/x-www-form-urlencoded'
            }
        }).success(function (data) {
      if( data.title == "Swim Meet" ) {
				$http({
            method : 'DELETE',
            url : '/gosched/v1/activity?id=' + data.id,
            headers : {
                'Content-Type' : 'application/x-www-form-urlencoded'
            }
        	}).success(function (data) {
      	if( data.message == "SUCCESS" ) {
					$scope.DeleteActivityTestData = setOutcome("Passed", data, "#00FF00");
				} else {
					$scope.DeleteActivityTestData = setOutcome("Failed", data, "#FF0000");
				}
			});
		}
    }).error(function (data) {
				$scope.DeleteActivityTestData = setOutcome("Failed", data, "#FF0000");
		});
  };

	$scope.DeleteEventTest = function() {
		$http({
            method : 'POST',
            url : '/gosched/v1/activity/events',
            data : 'title=Private Lesson&description=Flip Turns&activity_id=' + $scope.actid + '&start_time=Jan 2, 2013 at 3:00pm (MST)&end_time=Jan 2, 2013 at 4:00pm (MST)&max_attendees=8',
            headers : {
                'Content-Type' : 'application/x-www-form-urlencoded'
            }

        }).success(function (data) {
      if( data.title == "Private Lesson" ) {
				$http({
            method : 'DELETE',
            url : '/gosched/v1/activity/events?id=' + data.id,
            headers : {
                'Content-Type' : 'application/x-www-form-urlencoded'
            }
        	}).success(function (data) {
      	if( data.message == "SUCCESS" ) {
					$scope.DeleteEventTestData = setOutcome("Passed", data, "#00FF00");
				} else {
					$scope.DeleteEventTestData = setOutcome("Failed", data, "#FF0000");
				}
			});
		}
    }).error(function (data) {
				$scope.DeleteEventTestData = setOutcome("Failed", data, "#FF0000");
		});
  };

	$scope.GetActivityTest = function() {
		$http({
            method : 'GET',
            url : '/gosched/v1/activity?id=' + $scope.actid,
            headers : {
                'Content-Type' : 'application/x-www-form-urlencoded'
            }
        }).success(function (data) {
      if( data.title == "Swim Meet" ) {
				$scope.GetActivityTestData = setOutcome("Passed", data, "#00FF00");
			} else {
				$scope.GetActivityTestData = setOutcome("Failed", data, "#FF0000");
			}
    }).error(function (data) {
				$scope.GetActivityTestData = setOutcome("Failed", data, "#FF0000");
		});
  };

	$scope.GetEventTest = function() {
		$http({
            method : 'GET',
            url : '/gosched/v1/activity/events?id=' + $scope.eventid,
            headers : {
                'Content-Type' : 'application/x-www-form-urlencoded'
            }
        }).success(function (data) {
      if( data.title == "Private Lesson" ) {
				$scope.GetEventTestData = setOutcome("Passed", data, "#00FF00");
			} else {
				$scope.GetEventTestData = setOutcome("Failed", data, "#FF0000");
			}
    }).error(function (data) {
				$scope.GetEventTestData = setOutcome("Failed", data, "#FF0000");
		});
  };

  $scope.EventSearch1Test = function() {
    $http({
            method : 'GET',
            url : '/gosched/v1/activity/events/search?day=Jan 2, 2013&activity_id=' + $scope.actid,
            headers : {
                'Content-Type' : 'application/x-www-form-urlencoded'
            }
        }).success(function (data) {
      if( data instanceof Array ) {
        $scope.EventSearch1TestData = setOutcome("Passed", data, "#00FF00");
      } else {
        $scope.EventSearch1TestData = setOutcome("Failed", data, "#FF0000");
      }
    }).error(function (data) {
        $scope.EventSearch1TestData = setOutcome("Failed", data, "#FF0000");
    });
  };

	$scope.EventSearch2Test = function() {
		$http({
            method : 'GET',
            url : '/gosched/v1/activity/events/search?results=days&day=Jan 2, 2013&dayend=Jan 4, 2013&activity_id=' + $scope.actid,
            headers : {
                'Content-Type' : 'application/x-www-form-urlencoded'
            }
        }).success(function (data) {
      if( data instanceof Array ) {
				$scope.EventSearch2TestData = setOutcome("Passed", data, "#00FF00");
			} else {
				$scope.EventSearch2TestData = setOutcome("Failed", data, "#FF0000");
			}
    }).error(function (data) {
				$scope.EventSearch2TestData = setOutcome("Failed", data, "#FF0000");
		});
	};

	$scope.ClearData = function() {
		$http({
            method : 'GET',
            url : '/gosched/v1/activity/clear?owner=' + $scope.username,
            headers : {
                'Content-Type' : 'application/x-www-form-urlencoded'
            }
        }).success(function (data) {
      if( data.message == "SUCCESS" ) {
				$scope.initValues();
			} 
		});
	};

	$scope.initValues();
}

angular.module('gosched',[]).
  config(['$routeProvider', function($routeProvider) {
  $routeProvider.
      when('/tests', {templateUrl: '/partials/unittest.html', controller: UnitTestCtrl}).
			otherwise({redirectTo: '/tests'});
}]);
