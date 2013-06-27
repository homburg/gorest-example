(function () {
	var app = angular.module("myApp", ["ngResource"]);

	app.factory('users', function($resource) {
		return $resource('/service/users/:id');
	});

	app.factory("uptime", function($resource) {
		return $resource("/service/uptime");
	});

	app.directive("ngUserDetails", function () {
		return {
			restrict: "E",
			scope: {},
			template: '<div><h2>{{activeUser.Name}}</h2><p>{{activeUser.Email}}</p></div>'
		};
	});

	app.controller("UsersCtrl", ["$scope", "users", "uptime", function ($scope, users, uptime) {
		$scope.title = "Angular title!";
		var refreshUsers = function() {
			$scope.users = users.query()
		};

		var refreshUptime = function () {
			$scope.uptime = uptime.get()
		};

		$scope.AddUser = function () {
			users.save({
				Name: $scope.name,
				Email: $scope.email
			}, function() {
				refreshUsers()
			});
		};

		$scope.ViewUser = function(user) {
			$scope.activeUser = user;
		};

		refreshUsers();
	}]);
})();
