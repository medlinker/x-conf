$(function() {

    setPrjs();

    $(".createPrj").click(function() {
        $('.ui.small.modal.prj').modal('show');
    });

    $(".creatingPrj").click(function() {
        let prjName = $.trim($("#inputPrjName").val());
        if (prjName == "") {
            alert("请输入项目名");
            return;
        }
        $.post(
            "/x/conf/prj",
            {prjName: prjName},
            function(d) {
                if (d.code != 0) {
                    alert(d.msg);
                    return;
                }
                $("#inputPrjName").val("");
                setPrjs();
                $('.ui.small.modal.prj').modal('hide');
            },
            "json"
        );
    });

    $("#prjs").change(function() {
        $(".envMenu > a").removeClass("active");
    });

    $(".envMenu > a").click(function() {
        let env = $.trim($(this).text()).toLowerCase();
        let prj = $("#prjs").val();
        if (prj == "") {
            alert("请选择项目");
            return;
        }
        $(".envMenu > a").removeClass("active");
        $(this).addClass("active");
        $(".download").attr("href", "/x/conf/download?env="+env+"&prjName="+prj)
        confList(env, prj)
    });

    $(".configure").click(function() {
        let env = $.trim($('.envMenu > a[class="item active"]').text()).toLowerCase();
        let prj = $("#prjs").val();
        if (prj == "" || env == "") {
            alert("请选择项目和配置环境");
            return;
        }
        $('.ui.small.modal.conf').modal('show');
    });

    $(".createConf").click(function() {

        let env = $.trim($('.envMenu > a[class="item active"]').text()).toLowerCase();
        let prj = $("#prjs").val();
        if (prj == "" || env == "") {
            alert("请选择项目和配置环境");
            return;
        }

        let key = $.trim($("#inputKey").val());
        let value = $.trim($("#inputValue").val());

        if (key == "" || value == "") {
            alert("key and value required");
            return;
        }

        $.post(
            "/x/conf/configure",
            {env: env, prjName: prj, key: key, value: value},
            function(d) {
                if (d.code != 0) {
                    alert(d.msg);
                    return;
                }
                $("#inputKey").val("");
                $("#inputValue").val("");
                $('.ui.small.modal.conf').modal('hide');
                confList(env, prj)
            },
            "json"
        );
    });

    $(".upload").click(function() {
        let prj = $("#prjs").val();
        let env = $.trim($('.envMenu > a[class="item active"]').text()).toLowerCase();
        if (prj == "" || env == "") {
            alert("chose prj and env");
            return;
        }
        $("#inputUpPrj").val(prj);
        $("#inputUpEnv").val(env);
        $('.ui.small.modal.uploadM').modal('show');
    });

    $(".publish").click(function() {
        let prj = $("#prjs").val();
        let env = $.trim($('.envMenu > a[class="item active"]').text()).toLowerCase();
        if (prj == "" || env == "") {
            alert("chose prj and env");
            return;
        }
        $.post(
            "/x/conf/publish",
            {env: env, prjName: prj},
            function(d) {
                if (d.code != 0) {
                    alert(d.msg);
                    return;
                }
                alert("Publish Success");
            },
            "json"
        );
    });
    $(".heartCheck").click(function() {
        let prj = $("#prjs").val();
        let env = $.trim($('.envMenu > a[class="item active"]').text()).toLowerCase();
        if (prj == "" || env == "") {
            alert("chose prj and env");
            return;
        }
        $.post(
            "/x/conf/heartbeat",
            {env: env, prjName: prj},
            function(d) {
                if (d.code != 0) {
                    alert(d.msg);
                    return;
                }
                $(".heart-instance").empty();
                for (i in d.data) {
                    $(".heart-instance").append(`<tr><td>`+ d.data[i]
                        +`</td><td><i class="large green checkmark icon"></i></td></tr>`)
                }
                $('.ui.small.modal.heart').modal('show');
            },
            "json"
        );
    });
})

function setPrjs() {
    $.get(
        "/x/conf/prjs",
        function(d) {
            if (d.code != 0) {
                alert(d.msg);
                return;
            }
            $("#prjs").empty();
            $("#prjs").append('<option value="">select projects</option>');
            for (let prj in d.data) {
                $("#prjs").append('<option value="' + d.data[prj] +  '">' + d.data[prj] + '</option>"');
            }
        },
        "json"
    );
}

function notEmpty(key) {
    let val = $.trim($("#" + key).val());
    if (val == "") {
        return false;
    }
    return true;
}

function confList(env, prj) {
    $(".confList").empty();
    $.post(
        "/x/conf/configs",
        {env: env, prjName: prj},
        function(d) {
            if (d.code != 0) {
                alert(d.msg);
                return;
            }
            for (let c in d.data) {
                $(".confList").append(`
                    <tr>
                    <td>`+ c +`</td>
                    <td>`+ d.data[c] + `</td>
                    <td class="right aligned">
                        <i class="blue remove icon" onclick="delConf('` + prj + `','` + env + `','` + c + `')"></i>
                        <i class="blue file icon" onclick="updateConf('` + c + `','` + d.data[c] + `')"></i>
                    </td>
                    </tr>
                `);
            }
        },
        "json"
    );
}

function delConf(prj, env, key) {
    if (window.confirm("sure?")) {
        $.post(
            "/x/conf/del",
            {env: env, prjName: prj, key: key},
            function(d) {
                if (d.code != 0) {
                    alert(d.msg);
                    return;
                }
                confList(env, prj)
            },
            "json"
        );
    }
}

function updateConf(key, value) {
    $("#inputKey").val(key);
    $("#inputValue").val(value);
    $('.ui.small.modal.conf').modal('show');
}
