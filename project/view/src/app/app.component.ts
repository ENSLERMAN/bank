import {AfterContentInit, Component, DoCheck, OnDestroy, OnInit} from '@angular/core';
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
  ) { }

  async ngOnInit() {

  }

  ngDoCheck(): void {
    if (this.authGuard.isLoggedIn() == true) {
      this.visibility = true;
    } else if (this.authGuard.isLoggedIn() == false) {
      this.visibility = false;
    }

  }

  async getUser() {
    // @ts-ignore
    this.UserInfo = await this.http.getUserInfo();
  }

  logout() {
    console.log('logout');
    this.visibility = false;
    this.authService.logout();
    this.router.navigate(['/login']);
  }

}
