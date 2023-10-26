FROM fyneio/fyne-cross-images:linux

RUN apt update -y
RUN apt install -y libasound2-dev
RUN head -n -1 /usr/lib/x86_64-linux-gnu/pkgconfig/alsa.pc >/usr/lib/x86_64-linux-gnu/pkgconfig/tmp_alsa.pc && mv /usr/lib/x86_64-linux-gnu/pkgconfig/tmp_alsa.pc /usr/lib/x86_64-linux-gnu/pkgconfig/alsa.pc
RUN echo "Cflags: -I\${includedir}/alsa" >>/usr/lib/x86_64-linux-gnu/pkgconfig/alsa.pc