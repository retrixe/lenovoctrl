Name:           lenovoctrl
Version:        1.0.0
Release:        1%{?dist}
Summary:        Linux daemon and tools to control aspects of Lenovo IdeaPad/Legion devices.
BuildArch:      x86_64 aarch64
URL:            https://github.com/retrixe/lenovoctrl
Group:          System
Packager:       retrixe

License:        GPL-3.0-or-later
Source0:        %{expand:%%(pwd)}
#Source0:        {name}-{version}.tar.gz

# openSUSE calls it libappindicator3-1, but this works on both Fedora and openSUSE
Requires:       libappindicator-gtk3
%if 0%{?suse_version}
BuildRequires:  libappindicator3-devel
%else
BuildRequires:  libappindicator-gtk3-devel
%endif
BuildRequires:  golang
BuildRequires:  systemd
BuildRequires:  systemd-rpm-macros

%description
Linux daemon and tools to control aspects of Lenovo IdeaPad/Legion devices.
Currently supports toggling Battery Conservation Mode.

%prep
# we do a bit of tomfoolery
cd %{SOURCEURL0}
./scripts/build.sh
mkdir %{name}-%{version}
mv lenovoctrl %{name}-%{version}
tar --create --file %{name}-%{version}.tar.gz %{name}-%{version}
rm -r %{name}-%{version}
mv %{name}-%{version}.tar.gz %{_sourcedir}/%{name}

%setup -q

%install
rm -rf $RPM_BUILD_ROOT
# lenovoctrl binary
mkdir -p $RPM_BUILD_ROOT/%{_bindir}
cp %{name} $RPM_BUILD_ROOT/%{_bindir}
# lenovoctrl desktop file
mkdir -p $RPM_BUILD_ROOT/%{_datadir}/applications
cp %{SOURCEURL0}/snap/gui/lenovoctrl.desktop $RPM_BUILD_ROOT/%{_datadir}/applications
sed -i 's/\${SNAP}\/meta\/gui\/icon.png/lenovoctrl/g' $RPM_BUILD_ROOT/%{_datadir}/applications/lenovoctrl.desktop
# lenovoctrl icon
mkdir -p $RPM_BUILD_ROOT/%{_datadir}/icons/hicolor/symbolic/apps
cp %{SOURCEURL0}/snap/gui/icon.png $RPM_BUILD_ROOT/%{_datadir}/icons/hicolor/symbolic/apps/lenovoctrl.png
# lenovoctrl systemd service
mkdir -p $RPM_BUILD_ROOT/%{_unitdir}
mkdir -p $RPM_BUILD_ROOT/%{_presetdir}
cp %{SOURCEURL0}/scripts/packaging/lenovoctrl.service $RPM_BUILD_ROOT/%{_unitdir}
echo "enable lenovoctrl.service" > $RPM_BUILD_ROOT/%{_presetdir}/70-lenovoctrl.preset
# lenovoctrl d-bus policy
mkdir -p $RPM_BUILD_ROOT/%{_sysconfdir}/dbus-1/system.d
cp %{SOURCEURL0}/scripts/packaging/dbus-policy.conf $RPM_BUILD_ROOT/%{_sysconfdir}/dbus-1/system.d/com.retrixe.LenovoCtrl.v0.conf
# lenovoctrl license
cp %{SOURCEURL0}/LICENSE .

%clean
rm -rf $RPM_BUILD_ROOT

%post
%systemd_post lenovoctrl.service

%preun
%systemd_preun lenovoctrl.service

%postun
%systemd_postun_with_restart lenovoctrl.service

%files
%license LICENSE
%{_bindir}/%{name}
%{_datadir}/applications/lenovoctrl.desktop
# TODO: could be better
%{_datadir}/icons/hicolor/symbolic/apps/lenovoctrl.png
%{_unitdir}/lenovoctrl.service
%{_presetdir}/70-lenovoctrl.preset
%{_sysconfdir}/dbus-1/system.d/com.retrixe.LenovoCtrl.v0.conf

%changelog
# TODO: Not available on openSUSE!
%autochangelog
