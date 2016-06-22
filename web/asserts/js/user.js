$(function() {
    $(".submit").click(function() {
        let name = $.trim($("#name").val());
        let pass = $.trim($("#pass").val());
        if (name == "" || pass == "") {
            alert("username and password required!");
            return;
        }
        $.post(
            "http://127.0.0.1:8000/x/conf/login",
            {name: name, pass: pass},
            function(d) {
                if (d.code != 0) {
                    alert(d.msg);
                    return;
                }
                self.location='config.html'; 
            },
            "json"
        );
    });
})
