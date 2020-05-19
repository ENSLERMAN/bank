import { Component, OnInit } from '@angular/core';
import {FormControl, FormGroup, Validators} from "@angular/forms";
import {Router} from "@angular/router";
import {AuthService} from "../../services/auth.service";
import {HttpService} from "../../services/http.service";

@Component({
  selector: 'app-signin',
  templateUrl: './signin.component.html',
  styleUrls: ['./signin.component.scss']
})
export class SigninComponent implements OnInit {

  constructor(
    private router: Router,
    private authService: AuthService,
    private http: HttpService
  ) { }

  hide = true;
  loginForm: FormGroup;
  isMobile: boolean = false;
  returnUrl: string;

  ngOnInit(): void {

    if (document.documentElement.clientHeight > document.documentElement.clientWidth) {
      this.isMobile = true
    }

    this.loginForm = new FormGroup({
      email: new FormControl('', [Validators.required, Validators.minLength(6)]),
      password: new FormControl('', [Validators.required, Validators.minLength(6)])
    })
    this.returnUrl = '/main';
    this.authService.logout();
  }

  get f() { return this.loginForm.controls; }

  async login() {
    await this.http.authUser({
      "login": this.loginForm.value.email,
      "password": this.loginForm.value.password
    }).then((res) => {
      if (res["status"] == 200) {
        console.log("Login successful");
        //this.authService.authLogin(this.model);
        localStorage.setItem('isLoggedIn', "true");
        localStorage.setItem('token', this.f.email.value);
        this.router.navigate([this.returnUrl]);
      }
    });
  }

  rememberPassword() {
    alert("ЭТА ФУНКЦИЯ ЕЩЕ НЕ РЕАЛИЗОВАНА, НО КОГДА ОНА БУДЕТ СДЕЛАНА, ТЕБЕ ОБ ЭТОМ НИКТО НЕ СКАЖЕТ")
  }

  register() {
    this.router.navigate(['/register']);
  }

}
