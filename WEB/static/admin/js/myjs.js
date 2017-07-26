/**
 * Created by python on 17-6-6.
 */
$(function () {
    $(document).keydown(function (e) {
        if (e.keyCode === 13) {
            $("input:submit").filter(":eq(0)")[0].click()
        }
    });
    $("body").hide().fadeIn(1500);
    $("#showA").on('click', function () {
        $(this).animate({
            left: (200).toString() + 'px',
            opacity: 'hide'
        }, 1500);
        $.get('show/', '', function (data) {
            $("#showPage:eq(1)").html(data).hide().delay(1000).fadeIn(1500);
        });
        $("#showPage").load('show/').hide().delay(1000).fadeIn(1500);
    });
});