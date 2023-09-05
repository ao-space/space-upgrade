# Copyright (c) 2022 Institute of Software, Chinese Academy of Sciences (ISCAS)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

%global _bin_path /usr/local/bin
%global _service_path /usr/lib/systemd/system
%global debug_package %{nil}

Name:	 eulixspace-upgrade
Version: 1.0.4
Release: 13
Summary: upgrade for EulixOS box
License: Unlicense
URL:	 https://code.eulix.xyz/bp/box/system/aospace-upgrade
Source0: aospace-upgrade.tar.gz

BuildRequires: golang
BuildRequires: gcc >= 3.4.2
AutoReq: no
AutoProv: yes

Provides: eulixspace-upgrade = %{version}-%{release}
Provides: aospace-upgrade = %{version}-%{release}

ExclusiveArch: aarch64 x86_64
ExclusiveOS: Linux

%description
upgrade for EulixOS box

%prep
%setup -q -n %{name}-%{version} -c

%install
mkdir -p %{buildroot}%{_bin_path}
mkdir -p %{buildroot}%{_service_path}
cd aospace-upgrade
install -p -m 755 build/*-upgrade %{buildroot}%{_bin_path}
install -p -m 644 eulixspace-upgrade.service %{buildroot}%{_service_path}

%files
%defattr (-, root, root)
%{_bin_path}/*
%{_service_path}/*

%post
%systemd_post eulixspace-upgrade.service
systemctl enable eulixspace-upgrade.service
systemctl start eulixspace-upgrade.service

%preun
%systemd_preun eulixspace-upgrade.service

%postun
%systemd_postun_with_restart eulixspace-upgrade.service

%changelog
* Tue Jul 11 2023 Xuyang Zhang<zhangxuyang@iscas.ac.cn> - 1.0.4-13
- upgrade run in container

* Tue May 17 2023 Xuyang Zhang<zhangxuyang@iscas.ac.cn> - 1.0.4-6
- refactor eulixspace-upgrade

* Tue May 16 2023 Yafen Fang<yafen@iscas.ac.cn> - 1.0.4-5
- remove eulixspace-upgrade.service

* Wed Apr 20 2022 Jiayi Yin<jiayi@iscas.ac.cn> - 1.0.4-4
- Release 1.0.4-4

* Thu Feb 17 2022 Jiayi Yin<jiayi@iscas.ac.cn> - 1.0.1-3
- Release 1.0.1-3

* Fri Jan 21 2022 Jiayi Yin<jiayi@iscas.ac.cn> - 1.0.1-2
- Release 1.0.1-2

* Mon Jan 17 2022 Jiayi Yin<jiayi@iscas.ac.cn> - 1.0.1-1
- Release 1.0.1-1

* Thu Jan 06 2022 Jiayi Yin<jiayi@iscas.ac.cn> - 1.0.0-7
- Release 1.0.0-7

* Thu Jan 06 2022 Jiayi Yin<jiayi@iscas.ac.cn> - 1.0.0-6
- Release 1.0.0-6

* Wed Jan 05 2022 Jiayi Yin<jiayi@iscas.ac.cn> - 1.0.0-4
- Release 1.0.0-4

* Wed Jan 05 2022 Jiayi Yin<jiayi@iscas.ac.cn> - 1.0.0-2
- Release 1.0.0-2

* Wed Jan 05 2022 Jiayi Yin<jiayi@iscas.ac.cn> - 1.0.0-1
- Release 1.0.0-1

* Tue Jan 04 2022 Jiayi Yin<jiayi@iscas.ac.cn> - 0.7.6-2
- Release 0.7.6-2

* Tue Jan 04 2022 Jiayi Yin<jiayi@iscas.ac.cn> - 0.7.6-1
- Release 0.7.6-1

* Mon Jan 03 2022 Jiayi Yin<jiayi@iscas.ac.cn> - 0.7.4-9
- Release 0.7.4-9

* Mon Jan 03 2022 Jiayi Yin<jiayi@iscas.ac.cn> - 0.7.4-8
- Release 0.7.4-8

* Sat Jan 01 2022 Jiayi Yin<jiayi@iscas.ac.cn> - 0.7.4-7
- Release 0.7.4-7

* Sat Jan 01 2022 Jiayi Yin<jiayi@iscas.ac.cn> - 0.7.4-6
- Release 0.7.4-6

* Fri Dec 31 2021 Jiayi Yin<jiayi@iscas.ac.cn> - 0.7.4-5
- Release 0.7.4-5

* Fri Dec 31 2021 Jiayi Yin<jiayi@iscas.ac.cn> - 0.7.4-4
- Release 0.7.4-4

* Fri Dec 31 2021 Jiayi Yin<jiayi@iscas.ac.cn> - 0.7.4-3
- Release 0.7.4-3

* Thu Dec 30 2021 Jiayi Yin<jiayi@iscas.ac.cn> - 0.7.4-2
- Release 0.7.4-2

* Thu Dec 30 2021 Jiayi Yin<jiayi@iscas.ac.cn> - 0.7.4-1
- Release 0.7.4-1

* Thu Dec 30 2021 Yongxin Lei<Yongxin@iscas.ac.cn> - 0.0.1-2
- add recheck upgrade status

* Thu Dec 30 2021 Yafen Fang<yafen@iscas.ac.cn> - 0.0.1-1
- init package
