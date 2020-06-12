import {Component, DoCheck, OnInit} from '@angular/core';
import {Router} from "@angular/router";
import {AuthService} from "../../services/auth.service";
import {HttpService} from "../../services/http.service";
import {FormControl, FormGroup, Validators} from "@angular/forms";

@Component({
  selector: 'app-signup',
  templateUrl: './signup.component.html',
  styleUrls: ['./signup.component.scss']
})
export class SignupComponent implements OnInit, DoCheck {

  constructor(
      private router: Router,
      private authService: AuthService,
      private http: HttpService
  ) { }

  hide = true;
  loginForm: FormGroup;
  isMobile: boolean = false;
  returnUrl: string;
  error: boolean = false;

  ngOnInit(): void {
    if (document.documentElement.clientHeight > document.documentElement.clientWidth) {
      this.isMobile = true
    }
    this.loginForm = new FormGroup({
      login: new FormControl('', [Validators.required, Validators.minLength(6)]),
      password: new FormControl('', [Validators.required, Validators.minLength(6)]),
      name: new FormControl('', [Validators.required]),
      surname: new FormControl('', [Validators.required]),
      pat: new FormControl('', [Validators.required]),
      acceptTerms: new FormControl('', Validators.requiredTrue),
      sex: new FormControl(''),
      passport: new FormControl('', [Validators.required, Validators.pattern(RegExp("^[0-9]{10}$"))]),
    })
    this.returnUrl = '/main';
    this.authService.logout();
  }

  ngDoCheck(): void {
    if (document.documentElement.clientHeight > document.documentElement.clientWidth) {
      this.isMobile = true
    }
  }

  get f() { return this.loginForm.controls; }

  async register() {
    await this.http.registerUser({
      "login": this.loginForm.value.login,
      "password": this.loginForm.value.password,
      "surname": this.loginForm.value.surname,
      "name": this.loginForm.value.name,
      "patronymic": this.loginForm.value.pat,
      "passport": this.loginForm.value.passport,
    }).then(
        (res) => {
          if (res["status"] == 201) {
            console.log("Register successful");
            this.login()
          }
        },
        (err) => {
          console.log("jopa")
          this.error = true;
        }
    );
  }

  async login() {
    await this.http.authUser({
      "login": this.loginForm.value.email,
      "password": this.loginForm.value.password,
    }).then(
        (res) => {
          if (res["status"] == 200) {
            console.log("login successful");
            localStorage.setItem('isLoggedIn', "true");
            localStorage.setItem('token', this.f.email.value);
            this.router.navigate([this.returnUrl]);
          }
        },
    )
  }

}
