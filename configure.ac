# Release version
m4_define([package_version_major],[1])
m4_define([package_version_minor],[0])
m4_define([package_version_micro],[5])
 
# Initialize the process 
AC_INIT([],
        [package_version_major.package_version_minor.package_version_micro],
        …)
 
AC_SUBST([PACKAGE_VERSION_MAJOR],package_version_major)
AC_SUBST([PACKAGE_VERSION_MINOR],package_version_minor)
AC_SUBST([PACKAGE_VERSION_MICRO],package_version_micro)
 
# Output files
AC_CONFIG_FILES([
libuplinkc.pc:libuplinkc.pc.in
])

# init automake
AM_INIT_AUTOMAKE([1.11])
#configure and create "Makefile"
AC_CONFIG_FILES([Makefile])

AC_OUTPUT
