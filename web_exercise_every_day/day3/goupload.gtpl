<html>
    <head>
        <title>upload file</title>
    </head>
    <body>
        <form enctype="multipart/form-data" action="http://192.168.2.119:9090" method="post">
            <input type="file" name="uploadfile" />
            <input type="hidden" name="token" value="{{.}}"/>
            <input type="submit" value="upload" />
        </form>
    </body>
</html>
