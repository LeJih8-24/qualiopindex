import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

export interface Format {
  id: string;
  client: string;
  description: string;
  formation: string;
  date: string;
  formateur: string;
  financer: string;
  durer: number;
  satis: number;
  satisLater: number;
}
export interface Day {
  id: string;
  format_id: string;
  which: number;
  index: number;
}
export interface Present {
  id: string;
  day_id: string;
  present: string;
}

@Injectable({ providedIn: 'root' })
export class FormatService {
  private apiUrl = `${environment.apiUrl}/formats`;
  constructor(private http: HttpClient) {}

  list(): Observable<Format[]> {
    return this.http.get<Format[]>(this.apiUrl);
  }
  get(id: string): Observable<Format> {
    return this.http.get<Format>(`${this.apiUrl}/${id}`);
  }
  create(f: Format): Observable<any> {
    return this.http.post(this.apiUrl, f);
  }
  update(f: Format): Observable<any> {
    return this.http.put(`${this.apiUrl}/${f.id}`, f);
  }

  getDays(formatId: string) {
    return this.http.get<Day[]>(`${this.apiUrl}/${formatId}/days`);
  }
  addDay(d: Day) {
    return this.http.post(`${environment.apiUrl}/days`, d);
  }

  getPresents(formatId: string) {
    return this.http.get<Present[]>(`${this.apiUrl}/${formatId}/presents`);
  }

  addPresent(p: Present) {
    return this.http.post(`${environment.apiUrl}/presents`, p);
  }

  updatePresent(id: string, text: string) {
    return this.http.patch(`${environment.apiUrl}/presents/${id}`, { present: text });
  }
  
  remove(id: string) {
    return this.http.delete(`${this.apiUrl}/${id}`);
  }
}
