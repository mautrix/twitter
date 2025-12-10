package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import java.io.IOException;
import kotlin.Metadata;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u0000<\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\t\n\u0002\b\u0004\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\u0007\n\u0002\u0010\u000e\n\u0002\b\u0002\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\u0006\b\u0086\b\u0018\u0000 \u001d2\u00020\u0001:\u0002\u001e\u001dB\u001b\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0004\u001a\u0004\u0018\u00010\u0002¢\u0006\u0004\b\u0005\u0010\u0006J\u0017\u0010\n\u001a\u00020\t2\u0006\u0010\b\u001a\u00020\u0007H\u0016¢\u0006\u0004\b\n\u0010\u000bJ\u0012\u0010\f\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\f\u0010\rJ\u0012\u0010\u000e\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u000e\u0010\rJ(\u0010\u000f\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0004\u001a\u0004\u0018\u00010\u0002HÆ\u0001¢\u0006\u0004\b\u000f\u0010\u0010J\u0010\u0010\u0012\u001a\u00020\u0011HÖ\u0001¢\u0006\u0004\b\u0012\u0010\u0013J\u0010\u0010\u0015\u001a\u00020\u0014HÖ\u0001¢\u0006\u0004\b\u0015\u0010\u0016J\u001a\u0010\u001a\u001a\u00020\u00192\b\u0010\u0018\u001a\u0004\u0018\u00010\u0017HÖ\u0003¢\u0006\u0004\b\u001a\u0010\u001bR\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010\u001cR\u0016\u0010\u0004\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0004\u0010\u001c¨\u0006\u001f"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MediaDimensions;", "Lcom/bendb/thrifty/a;", "", "width", "height", "<init>", "(Ljava/lang/Long;Ljava/lang/Long;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/Long;", "component2", "copy", "(Ljava/lang/Long;Ljava/lang/Long;)Lcom/x/dmv2/thriftjava/MediaDimensions;", "", "toString", "()Ljava/lang/String;", "", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/Long;", "Companion", "MediaDimensionsAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class MediaDimensions implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final Long height;

    @JvmField
    @InterfaceC88465b
    public final Long width;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new MediaDimensionsAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MediaDimensions$MediaDimensionsAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/MediaDimensions;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/MediaDimensions;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/MediaDimensions;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class MediaDimensionsAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public MediaDimensions m85948read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Long lValueOf = null;
            Long lValueOf2 = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new MediaDimensions(lValueOf, lValueOf2);
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        C11272a.m14141a(protocol, b);
                    } else if (b == 10) {
                        lValueOf2 = Long.valueOf(protocol.mo14124H0());
                    } else {
                        C11272a.m14141a(protocol, b);
                    }
                } else if (b == 10) {
                    lValueOf = Long.valueOf(protocol.mo14124H0());
                } else {
                    C11272a.m14141a(protocol, b);
                }
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a MediaDimensions struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("MediaDimensions");
            if (struct.width != null) {
                protocol.mo14136v3("width", 1, (byte) 10);
                protocol.mo14121B3(struct.width.longValue());
            }
            if (struct.height != null) {
                protocol.mo14136v3("height", 2, (byte) 10);
                protocol.mo14121B3(struct.height.longValue());
            }
            protocol.mo14134i0();
        }
    }

    public MediaDimensions(@InterfaceC88465b Long l, @InterfaceC88465b Long l2) {
        this.width = l;
        this.height = l2;
    }

    public static /* synthetic */ MediaDimensions copy$default(MediaDimensions mediaDimensions, Long l, Long l2, int i, Object obj) {
        if ((i & 1) != 0) {
            l = mediaDimensions.width;
        }
        if ((i & 2) != 0) {
            l2 = mediaDimensions.height;
        }
        return mediaDimensions.copy(l, l2);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final Long getWidth() {
        return this.width;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final Long getHeight() {
        return this.height;
    }

    @InterfaceC88464a
    public final MediaDimensions copy(@InterfaceC88465b Long width, @InterfaceC88465b Long height) {
        return new MediaDimensions(width, height);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof MediaDimensions)) {
            return false;
        }
        MediaDimensions mediaDimensions = (MediaDimensions) other;
        return Intrinsics.m65267c(this.width, mediaDimensions.width) && Intrinsics.m65267c(this.height, mediaDimensions.height);
    }

    public int hashCode() {
        Long l = this.width;
        int iHashCode = (l == null ? 0 : l.hashCode()) * 31;
        Long l2 = this.height;
        return iHashCode + (l2 != null ? l2.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        return "MediaDimensions(width=" + this.width + ", height=" + this.height + Separators.RPAREN;
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}