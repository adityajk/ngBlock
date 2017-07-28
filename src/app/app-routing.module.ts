import { NgModule }              from '@angular/core';
import { RouterModule, Routes }  from '@angular/router';

import {PersonComponent} from "./user/person.component";
import {TableComponent} from "./table/table.component";
import {TypographyComponent} from "./typography/typography.component";
import {IconComponent} from "./icons/icon.component";
import {NotificationComponent} from "./notification/notification.component";
import {HeaderComponent} from "./header/header.component";
import {FooterComponent} from "./footer/footer.component";
import {SideNavComponent} from "./sidenav/sidenav.component";
import {HomeComponent} from "./home/home.component";

const appRoutes: Routes = [
    { path: "person", component: PersonComponent },
    { path: "table", component: TableComponent },
    { path: "librarybooks", component: TypographyComponent },
    { path: "bubblechart", component: IconComponent },
    { path: "notifications", component: NotificationComponent },
    { path: "home", component: HomeComponent },
    { path: '', component: HeaderComponent, outlet: 'header' },
    { path: '', component: FooterComponent, outlet: 'footer' },
    { path: '', component: SideNavComponent, outlet: 'navbar' }
];

@NgModule({
    imports: [
        RouterModule.forRoot(
            appRoutes,
            { enableTracing: true } // <-- debugging purposes only
        )
    ],
    exports: [
        RouterModule
    ]
})
export class AppRoutingModule { }