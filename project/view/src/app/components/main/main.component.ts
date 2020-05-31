import {Component, HostListener, OnInit} from '@angular/core';
import {HttpService} from "../../services/http.service";

@Component({
  selector: 'app-main',
  templateUrl: './main.component.html',
  styleUrls: ['./main.component.scss']
})
export class MainComponent implements OnInit {

  constructor(
      private http: HttpService
  ) { }

  UserInfo: {
    id: number,
    login: string,
    name: string,
    passport: string,
    patronymic: string,
    surname: string,
  };
  Hour: any;
  Privet: string;
  cols: number = 3;
  Bills: object;

  async ngOnInit() {
    this.Hour = new Date().getHours();
    // @ts-ignore
    this.UserInfo = await this.http.getUserInfo();
    await this.getUserBills(this.UserInfo.id)
    console.log(this.Bills)

    if (this.Hour >= 3 && this.Hour < 12) {
      this.Privet = "Доброе утро";
    } else if (this.Hour >= 12 && this.Hour < 18) {
      this.Privet = "Добрый день";
    } else if (this.Hour >= 18 && this.Hour < 24) {
      this.Privet = "Добрый вечер";
    } else if (this.Hour >= 0 && this.Hour < 3) {
      this.Privet = "Доброй ночи";
    }
      console.log(this.UserInfo)
    console.log(this.Hour)
  }

  async getUserBills(id) {
    await this.http.getUserBills().then((res) => {
      this.Bills = res["body"]
    });
    // @ts-ignore
    for (let item of this.Bills) {
      item.number = item.number % 10000
      switch (item.name) {
        case "Default bill": {
          item.img = "assets/icons/card.png";
          item.name = "Лицевой счет"
          break
        }
        case "Mastercard": {
          item.img = "assets/icons/mastercard.png"
          break
        }
        case "Visa": {
          item.img = "assets/icons/visa.png"
          break
        }
        case "Mir": {
          item.img = "assets/icons/mir.png"
          break
        }
      }
    }
  }


  @HostListener('window:resize', ['$event'])
  onResize(event) {
    this.adaptiveGrid()
  }

  adaptiveGrid() {
    if(document.documentElement.clientHeight > document.documentElement.clientWidth) {
      this.cols = 1;
    } else {
      this.cols = 3;
    }
  }

}
