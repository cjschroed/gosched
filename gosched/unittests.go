/*
;; Copyright (c) 2012, 2013 All Rights Reserved 
;; Author: Carl Schroeder
*/

package gosched 

import (
  "appengine"
	"appengine/user"
	"net/http"
  "html/template"
)

type myuser struct {
  Username string
  LogoutURL string
  Includes string
}

func UnitTestSection(w http.ResponseWriter, r *http.Request) {
  var m myuser
  c := appengine.NewContext(r)
  t,err := template.New("foo").Delims("[[", "]]").Parse(ShowUTTemplate)
  if err != nil {
    c.Infof("Error compiling template: %v", err.Error())
    http.Error(w, "Resource not found. (3)", http.StatusNotFound)
  } else {
    u := user.Current(c)
    if u == nil {
			m.Username = "Guest"
    } else {
      m.Username = u.Email
    }
    m.LogoutURL,_ = user.LogoutURL(c, "/")
    t.Execute(w, m)
  }
}

const ShowUTTemplate = 
`<html lang='en'>
<head>
<link href='/css/bootstrap.css' rel='stylesheet'>

<style type='text/css'>
html {
background: #ffffff;
}
</style>
<script>
  goscheduser = [[.Username]]
</script>
<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.7.1/jquery.min.js" type="text/javascript"></script>
<script src='https://www.google.com/jsapi' type='text/javascript'></script>
<script src='https://ajax.googleapis.com/ajax/libs/angularjs/1.0.4/angular.min.js'></script>
<script src='https://ajax.googleapis.com/ajax/libs/angularjs/1.0.4/angular-sanitize.min.js'></script>

<script src='/js/unittests.js' type='text/javascript'></script>
</head>
<body ng-app='gosched' ng-controller='UnitTestCtrl'>

<div class='navbar navbar-fixed-top'>
	<div class='navbar-inner'>
	<div class='container-fluid container-lemur-nav'>
      <a href='/' class='brand'>LemurWeb&nbsp;&nbsp;</a>
      <ul class='nav pull-right'>
        <li class='dropdown'>
						<a href='#' class='dropdown-toggle' data-toggle='dropdown'>
                <div id='usernameDropDown'>[[.Username]]</div>
                <div class='usernameCaret'>
                  <b class='caret'></b>
                </div>
              </a>
          <ul class='dropdown-menu'>
          	<li>
            	<a href='#' data-link='uvHelp'>Help</a>
            </li>
            <li>
              <a href='/members/welcome'>Getting Started</a>
            </li>
            <li class='divider'></li>
            <li>
            	<a href='/members/user-info'>Account Settings</a>
            </li>
            <li>
            	<a href='' data-link='installWebClipper'>Install Web Clipper</a>
            </li>
            <li class='divider'></li>
            <li>
              <a href='[[.LogoutURL]]'>Logout</a>
            </li>
          </ul>
        </li>
      </ul>
    </div>
  </div>
</div>

<div class='row' style='margin-top:50px;'>
</div>

<div class='span32'>
  <div class='lemurTitleBar' style='margin-top:40px;'>{{message}}</div>
	<div class='row'>
		<div class='span24' ng-view> </div>
	</div>
</div>


</body>
</html>`

