import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import {SigninComponent} from "./components/signin/signin.component";
import {MainComponent} from "./components/main/main.component";


const routes: Routes = [
  { path: "login", component: SigninComponent },
  { path: "", component: MainComponent }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
