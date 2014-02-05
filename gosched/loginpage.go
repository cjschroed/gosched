package gosched 

const LoginPageTemplate =
`<html lang="en">
  <head>
<meta charset="utf-8">

    <title>Go Scheduling API</title>
<meta name="description" content="">
<meta name="author" content="">
<meta http-equiv="x-ua-compatible" content="IE=8"><meta http-equiv="x-ua-compatible" content="IE=9">
<!-- Le styles -->
<link href="/css/bootstrap.css" rel="stylesheet">

    <style type="text/css">

html {
background: #ffffff;
}
    </style>
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/1.7.1/jquery.min.js" type="text/javascript"></script>
    <script src="https://www.google.com/jsapi" type="text/javascript"></script>
		<script src='https://ajax.googleapis.com/ajax/libs/angularjs/1.0.4/angular.min.js'></script>
		<script src='https://ajax.googleapis.com/ajax/libs/angularjs/1.0.4/angular-sanitize.min.js'></script>
	 	<script src="/js/loginpage.js" type="text/javascript"></script>
  </head>
  <body ng-app='gosched'>
    <div class="navbar navbar-fixed-top">
      <div class="navbar-inner">
        <div class="container-fluid">
          <a href="/" class="brand">Go Scheduling API&nbsp;&nbsp;
          </a>
        </div>
      </div>
    </div>
    <div class="container-fluid">
      <div class="row" style="height:40px;">
      </div>
			<div>
      	<h3>Press a button to sign in with the designated provider:</h3>
				<br>
      </div>
			<div>
      	<a class="btn btn-openid" href="[[.GoogleLogin]]" alt="Login with your Google account.">
        	<img style="height:40px; padding-top:3px;" src="/images/google_logo.png" id="google_button"/>
        </a>
				<br><br>
        <a class="btn btn-openid" href="[[.YahooLogin]]" alt="Login with your Yahoo account.">
        	<img alt="yahoo login" style="height:26px;" src="/images/yahoo_logo.png" id="yahoo_button"/>
        </a>
       	<br/>
      </div>
    </div>
	<div ng-view></div>
  </body>
</html>`
