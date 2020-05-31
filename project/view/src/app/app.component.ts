import {Component, DoCheck} from '@angular/core';
import {Router} from "@angular/router";
import {AuthService} from "./services/auth.service";
import {AuthGuard} from "./guards/auth.guard";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements DoCheck {
  title = 'view';
  visibility: boolean = false;

  constructor(
      private router: Router,
      private authService: AuthService,
      private authGuard: AuthGuard,
  ) { }

  ngDoCheck(): void {
    if (this.authGuard.isLoggedIn() == true) {
      this.visibility = true;
    } else if (this.authGuard.isLoggedIn() == false) {
      this.visibility = false;
    }
  }

  logout() {
    console.log('logout');
    this.visibility = false;
    this.authService.logout();
    this.router.navigate(['/login']);
  }
}
