Name: sqlchk
Version: v0.1.1
Release: 1%{?dist}
Summary: Checks SQL 

License: BSD
URL: http://unherd.info/info
Source0: %{name}-%{version}.tar.gz


%description
A Checks SQL

%global debug_package %{nil}

%prep
%autosetup

%build
make build

%install
install -m 0755 -d $RPM_BUILD_ROOT/usr/local/bin
install -m 0755 sqlchk $RPM_BUILD_ROOT/usr/local/bin/sqlchk

%files
%defattr(-,root,root,-)
/usr/local/bin/sqlchk

%clean
rm -rf %{buildroot}

%changelog
* Sat Dec 16 2023 Joe Siwiak <kf4jas@gmail.com> Refactored spec file to work better
* Tue Nov 29 2022 Joe Siwiak <kf4jas@gmail.com>
