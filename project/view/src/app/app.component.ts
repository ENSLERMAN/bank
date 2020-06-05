import {Component, DoCheck, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from "@angular/router";
import {AuthService} from "./services/auth.service";
import {AuthGuard} from "./guards/auth.guard";
import {HttpService} from "./services/http.service";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit, DoCheck {

  title = 'view';
  visibility: boolean = false;
  visibilityBack: boolean = false;
  Text: string;

  UserInfo: {
    id: number,
    login: string,
    name: string,
    passport: string,
    patronymic: string,
    surname: string,
  };
  Hour: any;

  constructor(
      private router: Router,
      private authService: AuthService,
      private authGuard: AuthGuard,
      private http: HttpService,
      private route: ActivatedRoute
  ) { }

  async ngOnInit() {
    this.Hour = new Date().getHours();
    // @ts-ignore
    this.UserInfo = await this.http.getUserInfo();

    if (this.router.url === "/main") {
      this.visibilityBack = false;
    } else {
      this.visibilityBack = true;
    }
  }

  ngDoCheck(): void {
    if (this.authGuard.isLoggedIn() == true) {
      this.visibility = true;
    } else if (this.authGuard.isLoggedIn() == false) {
      this.visibility = false;
    }

    if (this.router.url === "/main") {
      this.visibilityBack = false;
      if (this.Hour >= 3 && this.Hour < 12) {
        this.Text = "Доброе утро, " + this.UserInfo?.name + " " + this.UserInfo?.patronymic + "!";
      } else if (this.Hour >= 12 && this.Hour < 18) {
        this.Text = "Добрый день, " + this.UserInfo?.name + " " + this.UserInfo?.patronymic + "!";
      } else if (this.Hour >= 18 && this.Hour < 24) {
        this.Text = "Добрый вечер, " + this.UserInfo?.name + " " + this.UserInfo?.patronymic + "!";
      } else if (this.Hour >= 0 && this.Hour < 3) {
        this.Text = "Доброй ночи, " + this.UserInfo?.name + " " + this.UserInfo?.patronymic + "!";
      }
    } else if (this.router.url === "/history") {
      this.Text = "История переводов"
      this.visibilityBack = true;
    } else if (this.router.url === "/send_money") {
      this.Text = "Перевод денег"
      this.visibilityBack = true;
    }

  }

  logout() {
    console.log('logout');
    this.visibility = false;
    this.authService.logout();
    this.router.navigate(['/login']);
  }
}
