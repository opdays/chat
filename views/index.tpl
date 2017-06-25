<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<title>Bootstrap 实例 - 超大屏幕（Jumbotron）</title>
    {{ assets_css "static/css/bootstrap.min.css"}}
    {{ assets_js "static/js/jquery.js"}}
    {{ assets_js "static/js/bootstrap.min.js"}}
</head>
<body>

<div class="container">
	<div class="jumbotron">
		<h1>{{ .IP}}</h1>
		<p>{{ .Username}}</p>
		<p><a class="btn btn-primary btn-lg" role="button">
			进入</a>
		</p>
	</div>
    <div class="row">
        <div class="col-md-12 col-lg-12">
            <form action="/" method="post">
                <textarea name="shell"></textarea>
                <input type="submit" value="确定">
            </form>
            <pre>
                {{.Result}}
            </pre>
        </div>

    </div>

</div>

</body>
</html>