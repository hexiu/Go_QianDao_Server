{{define "navbar"}}
<div class="container">
    <a class="navbar-brand" href="/">智邮普创</a>


    <ul class="nav navbar-nav navbar-right">
        {{if .IsLogin}}
        <li><a href="/userlist">用户列表</a></li>

        <li><a href="/login?exit=true">退出</a></li>
        {{else}}
        <li><a href="/login">管理员登录</a></li>
        {{end}}
    </ul>
</div>
{{end}}