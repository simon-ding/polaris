import 'package:flutter/material.dart';
import 'package:flutter/scheduler.dart' show timeDilation;
import 'package:flutter_login/flutter_login.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:quiver/strings.dart';
import 'package:ui/providers/login.dart';
import 'package:ui/weclome.dart';

class LoginScreen extends ConsumerWidget {
  static const route = '/login';

  const LoginScreen({super.key});

  Duration get loginTime => Duration(milliseconds: timeDilation.ceil() * 2250);


  Future<String?> _recoverPassword(String name) {
    return Future.delayed(loginTime).then((_) {
      return null;
    });
  }


  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return FlutterLogin(
      title: 'Polaris',
      onLogin: (data) {
        ref.read(authSettingProvider.notifier).login(data.name, data.password);
      },
      onSubmitAnimationCompleted: () {
        context.go(WelcomePage.route);
      },
      onRecoverPassword: _recoverPassword,
      userValidator: (value) => isBlank(value)? "不能为空":null,
      userType: LoginUserType.name,
      hideForgotPasswordButton: true,
      messages: LoginMessages(
        userHint: '用户名',
        passwordHint: '密码',
        loginButton: '登录',
      ),
    );
  }
}

class IntroWidget extends StatelessWidget {
  const IntroWidget({super.key});

  @override
  Widget build(BuildContext context) {
    return const Column(
      children: [
        Text.rich(
          TextSpan(
            children: [
              TextSpan(
                text: "You are trying to login/sign up on server hosted on ",
              ),
              TextSpan(
                text: "example.com",
                style: TextStyle(fontWeight: FontWeight.bold),
              ),
            ],
          ),
          textAlign: TextAlign.justify,
        ),
        Row(
          children: <Widget>[
            Expanded(child: Divider()),
            Padding(
              padding: EdgeInsets.all(8.0),
              child: Text("Authenticate"),
            ),
            Expanded(child: Divider()),
          ],
        ),
      ],
    );
  }
}
