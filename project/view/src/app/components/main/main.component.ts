import {Component, HostListener, Inject, OnDestroy, OnInit} from '@angular/core';
import { HttpService } from "../../services/http.service";
import {MAT_DIALOG_DATA, MatDialog, MatDialogRef} from "@angular/material/dialog";
import {Router} from "@angular/router";

export interface DialogData {
  TypeBill: number;
}

@Component({
  selector: 'app-main',
  templateUrl: './main.component.html',
  styleUrls: ['./main.component.scss']
})
export class MainComponent implements OnInit {

  constructor(
      private http: HttpService,
      public dialog: MatDialog,
      private router: Router,
  ) { }

  UserInfo: {
    id: number,
    login: string,
    name: string,
    passport: string,
    patronymic: string,
    surname: string,
  };
  cols: number = 3;
  Bills: object;
  Text: string;
  Hour: any;

  TypeBill: number;
  visibility: boolean = true;

  ngOnInit() {
    this.Hour = new Date().getHours();
    this.getUser().then(
        (async r => {
          console.log("данные пользователя получены")
          await this.getUserBills(this.UserInfo.id)
          if (this.Hour >= 3 && this.Hour < 12) {
            this.Text = "Доброе утро, " + this.UserInfo?.name + " " + this.UserInfo?.patronymic + "!";
          } else if (this.Hour >= 12 && this.Hour < 18) {
            this.Text = "Добрый день, " + this.UserInfo?.name + " " + this.UserInfo?.patronymic + "!";
          } else if (this.Hour >= 18 && this.Hour < 24) {
            this.Text = "Добрый вечер, " + this.UserInfo?.name + " " + this.UserInfo?.patronymic + "!";
          } else if (this.Hour >= 0 && this.Hour < 3) {
            this.Text = "Доброй ночи, " + this.UserInfo?.name + " " + this.UserInfo?.patronymic + "!";
          }
        }),
        (err => {
          console.log("неудалось получить данные пользователя")
        })
    );
  }

  async getUserBills(id: number) {
    await this.http.getUserBills().then((res) => {
      this.Bills = res["body"]
    }).finally(() => {
      this.visibility = false;
    });
    // @ts-ignore
    for (let item of this.Bills) {
      item.number = item.number % 10000
      switch (item.type) {
        case 1: {
          item.img = "assets/icons/card.png";
          item.name = "Лицевой счет"
          break
        }
        case 2: {
          item.img = "assets/icons/mastercard.png"
          item.name = "Mastercard"
          break
        }
        case 3: {
          item.img = "assets/icons/visa.png"
          item.name = "Visa"
          break
        }
        case 4: {
          item.img = "assets/icons/mir.png"
          item.name = "МИР"
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

  async getUser() {
    // @ts-ignore
    this.UserInfo = await this.http.getUserInfo();
  }

  openDialog(): void {
    const dialogRef = this.dialog.open(DialogWindow, {
      width: 'auto',
      minWidth: '300px',
      data: { value: this.TypeBill }
    });

    dialogRef.afterClosed().subscribe(result => {
      console.log('The dialog was closed');
      if (result !== undefined) {
        this.TypeBill = result;
        this.createBill(this.TypeBill).then();
        console.log(this.TypeBill);
      }
    });
  }

  async createBill(type) {
    await this.http.createBill({
      "type": parseInt(type),
    }).then(()=> {
      this.getUserBills(this.UserInfo.id)
    })
  }

}


@Component({
  selector: 'dialog-window',
  templateUrl: 'dialog-window.html',
  styles: [
      `.mobile {width: 100%; margin-left: 0 !important; margin-bottom: 0.5rem; }`
  ]
})
export class DialogWindow {

  isMobile: boolean = false;

  constructor(
      public dialogRef: MatDialogRef<DialogWindow>,
      @Inject(MAT_DIALOG_DATA) public data: DialogData
  ) {
    if (document.documentElement.clientHeight > document.documentElement.clientWidth) {
      this.isMobile = true
    }
  }

  onNoClick(): void {
    this.dialogRef.close();
  }

}