# $NetBSD: Makefile,v 1.16 2018/12/13 19:52:08 adam Exp $
#

PKGREVISION=	11
CATEGORIES=	math

MAINTAINER=	pkgsrc-users@NetBSD.org
HOMEPAGE=	${R_HOMEPAGE_BASE}/ncdf/
COMMENT=	High-level R interface to Unidata's netCDF data files
LICENSE=	gnu-gpl-v2

R_PKGNAME=	ncdf
R_PKGVER=	1.6.6

CONFIGURE_ARGS+=	--with-netcdf_incdir=${BUILDLINK_PREFIX.netcdf}/include
CONFIGURE_ARGS+=	--with-netcdf_libdir=${BUILDLINK_PREFIX.netcdf}/lib

.include "../../math/R/Makefile.extension"
.include "../../devel/netcdf/buildlink3.mk"
.include "../../mk/bsd.pkg.mk"
