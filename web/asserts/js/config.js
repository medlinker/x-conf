$(function() {

    setPrjs();

    $(".createPrj").click(function() {
        $('.ui.small.modal').modal('show');
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
                $('.ui.small.modal').modal('hide');
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
