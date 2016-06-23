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
            "http://127.0.0.1:8000/x/conf/prj",
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
        $(".envMenu > a").removeClass("active");
        $(this).addClass("active");
        let env = $.trim($(this).text()).toLowerCase();
        let prj = $("#prjs").val();
        if (prj != "") {

        }
    });

    $(".configure").click(function() {

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
            "http://127.0.0.1:8000/x/conf/configure",
            {env: env, prjName: prj, key: key, value: value},
            function(d) {
                if (d.code != 0) {
                    alert(d.msg);
                    return;
                }
                $("#inputKey").val("");
                $('.ui.small.modal.conf').modal('hide');
            },
            "json"
        );
    });
})

function setPrjs() {
    $.get(
        "http://127.0.0.1:8000/x/conf/prjs",
        function(d) {
            if (d.code != 0) {
                alert(d.msg);
                return;
            }
            $("#prjs").empty();
            $("#prjs").append('<option value="">选择配置项目</option>');
            for (let prj in d.data) {
                $("#prjs").append('<option value="' + d.data[prj] +  '">' + d.data[prj] + '</option>"');
            }
        },
        "json"
    );
}
