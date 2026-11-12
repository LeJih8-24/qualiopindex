import { Component, OnInit, signal } from '@angular/core';
import { Session } from './session/session';
import { Format, FormatService } from './format.service';

@Component({
  imports: [Session],
  selector: 'app-home',
  styleUrl: './home.scss',
  templateUrl: './home.html',
})
export class Home implements OnInit {
  formats = signal<Format[]>([]);
  selectedId = signal<string | null>(null);

  constructor(private formatService: FormatService) {}

  ngOnInit() {
    this.refresh();
  }

  refresh() {
    this.formatService.list().subscribe((f) => this.formats.set(f));
  }

  newSession() {
    const f: Format = {
      id: crypto.randomUUID(),
      client: '',
      description: '',
      formation: '',
      date: '',
      formateur: '',
      financer: '',
      durer: 0,
      satis: 0,
      satisLater: 0,
    };
    this.formatService.create(f).subscribe(() => this.selectedId.set(f.id));
  }

  open(id: string) {
    this.selectedId.set(id);
  }

  close() {
    this.selectedId.set(null);
    this.refresh();
  }

  confirmingId = signal<string | null>(null);

  askDelete(id: string) {
    this.confirmingId.set(id);
  }
  cancelDelete() {
    this.confirmingId.set(null);
  }

  confirmDelete(id: string) {
    this.formatService.remove(id).subscribe({
      next: () => {
        this.confirmingId.set(null);
        this.refresh();
      },
      error: (err) => console.error('delete failed', err),
    });
  }
}
